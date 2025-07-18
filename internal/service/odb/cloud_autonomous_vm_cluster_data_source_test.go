//Copyright © 2025, Oracle and/or its affiliates. All rights reserved.

package odb_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/odb"
	odbtypes "github.com/aws/aws-sdk-go-v2/service/odb/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tfodb "github.com/hashicorp/terraform-provider-aws/internal/service/odb"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type autonomousVMClusterDSTest struct {
	exaInfraDisplayNamePrefix            string
	odbNetDisplayNamePrefix              string
	autonomousVmClusterDisplayNamePrefix string
}

var autonomousVMClusterDSTestEntity = autonomousVMClusterDSTest{
	exaInfraDisplayNamePrefix:            "Ofake-exa",
	odbNetDisplayNamePrefix:              "odb-net",
	autonomousVmClusterDisplayNamePrefix: "Ofake-avmc",
}

func TestAccODBCloudAutonomousVmClusterDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)

	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
	var avmc1 odbtypes.CloudAutonomousVmCluster
	dataSourceName := "data.aws_odb_cloud_autonomous_vm_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			autonomousVMClusterDSTestEntity.testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ODBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             autonomousVMClusterDSTestEntity.testAccCheckCloudAutonomousVmClusterDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: autonomousVMClusterDSTestEntity.avmcBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "display_name", *avmc1.DisplayName),
				),
			},
		},
	})
}
func (autonomousVMClusterDSTest) checkCloudAutonomousVmClusterExists(ctx context.Context, name string, cloudAutonomousVMCluster *odbtypes.CloudAutonomousVmCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameCloudAutonomousVmCluster, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameCloudAutonomousVmCluster, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)
		resp, err := autonomousVMClusterDSTestEntity.findAVMC(ctx, conn, rs.Primary.ID)
		if err != nil {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameCloudAutonomousVmCluster, rs.Primary.ID, err)
		}

		*cloudAutonomousVMCluster = *resp

		return nil
	}
}

func (autonomousVMClusterDSTest) findAVMC(ctx context.Context, conn *odb.Client, id string) (*odbtypes.CloudAutonomousVmCluster, error) {
	input := odb.GetCloudAutonomousVmClusterInput{
		CloudAutonomousVmClusterId: aws.String(id),
	}
	out, err := conn.GetCloudAutonomousVmCluster(ctx, &input)
	if err != nil {
		if errs.IsA[*odbtypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: &input,
			}
		}
		return nil, err
	}

	if out == nil || out.CloudAutonomousVmCluster == nil {
		return nil, tfresource.NewEmptyResultError(&input)
	}

	return out.CloudAutonomousVmCluster, nil
}

func (autonomousVMClusterDSTest) testAccCheckCloudAutonomousVmClusterDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_odb_cloud_autonomous_vm_cluster" {
				continue
			}

			_, err := tfodb.FindCloudAutonomousVmClusterByID(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				return nil
			}
			if err != nil {
				return create.Error(names.ODB, create.ErrActionCheckingDestroyed, tfodb.ResNameCloudAutonomousVmCluster, rs.Primary.ID, err)
			}

			return create.Error(names.ODB, create.ErrActionCheckingDestroyed, tfodb.ResNameCloudAutonomousVmCluster, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}
func (autonomousVMClusterDSTest) testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)

	input := &odb.ListCloudAutonomousVmClustersInput{}

	_, err := conn.ListCloudAutonomousVmClusters(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func (autonomousVMClusterDSTest) avmcBasic() string {

	exaInfraDisplayName := sdkacctest.RandomWithPrefix(autonomousVMClusterDSTestEntity.exaInfraDisplayNamePrefix)
	odbNetworkDisplayName := sdkacctest.RandomWithPrefix(autonomousVMClusterDSTestEntity.odbNetDisplayNamePrefix)
	avmcDisplayName := sdkacctest.RandomWithPrefix(autonomousVMClusterDSTestEntity.autonomousVmClusterDisplayNamePrefix)

	exaInfraRes := autonomousVMClusterDSTestEntity.exaInfra(exaInfraDisplayName)
	odbNetRes := autonomousVMClusterDSTestEntity.odbNet(odbNetworkDisplayName)
	res := fmt.Sprintf(`
%s

%s

data "aws_odb_db_servers_list" "test" {
  cloud_exadata_infrastructure_id = aws_odb_cloud_exadata_infrastructure.test.id
}

resource "aws_odb_cloud_autonomous_vm_cluster" "test" {
    cloud_exadata_infrastructure_id         = aws_odb_cloud_exadata_infrastructure.test.id
  	odb_network_id                          =aws_odb_network.test.id
	display_name             				= %[3]q
  	autonomous_data_storage_size_in_tbs     = 5
  	memory_per_oracle_compute_unit_in_gbs   = 2
  	total_container_databases               = 1
  	cpu_core_count_per_node                 = 40
    license_model                                = "LICENSE_INCLUDED"
    db_servers								   = [ for db_server in data.aws_odb_db_servers_list.test.db_servers : db_server.id]
    scan_listener_port_tls = 8561
    scan_listener_port_non_tls = 1024
    maintenance_window = {
		 preference = "NO_PREFERENCE"
		 days_of_week =	[]
         hours_of_day =	[]
         months = []
         weeks_of_month =[]
         lead_time_in_weeks = 0
    }

}


data "aws_odb_cloud_autonomous_vm_cluster" "test" {
  id             = aws_odb_cloud_autonomous_vm_cluster.test.id

}
`, exaInfraRes, odbNetRes, avmcDisplayName)

	return res
}

func (autonomousVMClusterDSTest) odbNet(odbNetName string) string {
	networkRes := fmt.Sprintf(`


resource "aws_odb_network" "test" {
  display_name          = %[1]q
  availability_zone_id = "use1-az6"
  client_subnet_cidr   = "10.2.0.0/24"
  backup_subnet_cidr   = "10.2.1.0/24"
  s3_access = "DISABLED"
  zero_etl_access = "DISABLED"
}

`, odbNetName)
	return networkRes
}

func (autonomousVMClusterDSTest) exaInfra(exaInfraName string) string {
	exaInfraRes := fmt.Sprintf(`


resource "aws_odb_cloud_exadata_infrastructure" "test" {
  display_name          = %[1]q
  shape             	= "Exadata.X9M"
  storage_count      	= 3
  compute_count         = 2
  availability_zone_id 	= "use1-az6"
  customer_contacts_to_send_to_oci = ["abc@example.com"]
maintenance_window = {
  		custom_action_timeout_in_mins = 16
		days_of_week =	[]
        hours_of_day =	[]
        is_custom_action_timeout_enabled = true
        lead_time_in_weeks = 0
        months = []
        patching_mode = "ROLLING"
        preference = "NO_PREFERENCE"
		weeks_of_month =[]
  }
}

`, exaInfraName)
	return exaInfraRes
}
