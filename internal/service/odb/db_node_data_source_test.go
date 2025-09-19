// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package odb_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/odb"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfodb "github.com/hashicorp/terraform-provider-aws/internal/service/odb"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type testDbNodeDataSourceTest struct {
	exaDisplayNamePrefix             string
	oracleDBNetworkDisplayNamePrefix string
	vmClusterDisplayNamePrefix       string
}

var dbNodeDataSourceTestEntity = testDbNodeDataSourceTest{
	exaDisplayNamePrefix:             "Ofake-exa",
	oracleDBNetworkDisplayNamePrefix: "odb-net",
	vmClusterDisplayNamePrefix:       "Ofake-vmc",
}

// Acceptance test access AWS and cost money to run.
func TestAccODBDbNodeDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
	var dbServer odb.GetDbServerOutput
	exaInfraDisplayName := sdkacctest.RandomWithPrefix(dbServersListDataSourceTests.displayNamePrefix)
	dataSourceName := "data.aws_odb_db_server.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ODBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             dbServerDataSourceTestEntity.testAccCheckDbServersDestroyed(ctx),
		Steps: []resource.TestStep{
			{
				Config: dbServerDataSourceTestEntity.basic(dbServerDataSourceTestEntity.exaInfra(exaInfraDisplayName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					dbServerDataSourceTestEntity.testAccCheckDbServerExists(ctx, dataSourceName, &dbServer),
				),
			},
		},
	})
}

func (testDbNodeDataSourceTest) testAccCheckDbNodeExists(ctx context.Context, name string, output *odb.GetDbServerOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.DSNameDbServer, name, errors.New("not found"))
		}
		conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)
		var dbServerId = rs.Primary.ID
		var attributes = rs.Primary.Attributes
		exaId := attributes["exadata_infrastructure_id"]
		resp, err := dbServerDataSourceTestEntity.findDbServer(ctx, conn, &dbServerId, &exaId)
		if err != nil {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.DSNameDbServer, rs.Primary.ID, err)
		}
		*output = *resp
		return nil
	}
}

func (testDbNodeDataSourceTest) testAccCheckDbNodeDestroyed(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_odb_cloud_exadata_infrastructure" {
				continue
			}
			_, err := dbServerDataSourceTestEntity.findExaInfra(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				return nil
			}
			if err != nil {
				return create.Error(names.ODB, create.ErrActionCheckingDestroyed, tfodb.DSNameDbServer, rs.Primary.ID, err)
			}
			return create.Error(names.ODB, create.ErrActionCheckingDestroyed, tfodb.DSNameDbServer, rs.Primary.ID, errors.New("not destroyed"))
		}
		return nil
	}
}

func (testDbNodeDataSourceTest) findDbNode(ctx context.Context, conn *odb.Client, dbServerId *string, exaInfraId *string) (*odb.GetDbServerOutput, error) {
	inputWithExaId := &odb.GetDbServerInput{
		DbServerId:                   dbServerId,
		CloudExadataInfrastructureId: exaInfraId,
	}
	output, err := conn.GetDbServer(ctx, inputWithExaId)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (testDbNodeDataSourceTest) basic() string {

	vmClusterConfig := dbNodeDataSourceTestEntity.vmClusterBasicConfig()

	return fmt.Sprintf(`
%s

data "aws_odb_db_nodes_list" "test" {
  cloud_vm_cluster_id = aws_odb_cloud_vm_cluster_id.test.id
}

data "aws_odb_db_node" "test" {
  id = data.aws_odb_db_nodes_list.test.db_nodes[0].id
  cloud_vm_cluster_id = cloud_vm_cluster_id.test.id
}



`, vmClusterConfig)
}

func (testDbNodeDataSourceTest) vmClusterBasicConfig() string {

	exaInfraDisplayName := sdkacctest.RandomWithPrefix(dbNodeDataSourceTestEntity.exaDisplayNamePrefix)
	oracleDBNetDisplayName := sdkacctest.RandomWithPrefix(dbNodeDataSourceTestEntity.oracleDBNetworkDisplayNamePrefix)
	vmcDsplayName := sdkacctest.RandomWithPrefix(dbNodeDataSourceTestEntity.vmClusterDisplayNamePrefix)
	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		panic(err)
	}
	dsTfCodeVmCluster := fmt.Sprintf(`


resource "aws_odb_network" "test" {
  display_name         = %[1]q
  availability_zone_id = "use1-az6"
  client_subnet_cidr   = "10.2.0.0/24"
  backup_subnet_cidr   = "10.2.1.0/24"
  s3_access            = "DISABLED"
  zero_etl_access      = "DISABLED"
}

resource "aws_odb_cloud_exadata_infrastructure" "test" {
  display_name         = %[1]q
  shape                = "Exadata.X9M"
  storage_count        = 3
  compute_count        = 2
  availability_zone_id = "use1-az6"
  maintenance_window {
    custom_action_timeout_in_mins    = 16
    is_custom_action_timeout_enabled = true
    patching_mode                    = "ROLLING"
    preference                       = "NO_PREFERENCE"
  }
}

data "aws_odb_db_servers_list" "test" {
  cloud_exadata_infrastructure_id = aws_odb_cloud_exadata_infrastructure.test.id
}

resource "aws_odb_cloud_vm_cluster" "test" {
  display_name                    = %[3]q
  cloud_exadata_infrastructure_id = aws_odb_cloud_exadata_infrastructure.test.id
  cpu_core_count                  = 6
  gi_version                      = "23.0.0.0"
  hostname_prefix                 = "apollo12"
  ssh_public_keys                 = ["%[4]s"]
  odb_network_id                  = aws_odb_network.test.id
  is_local_backup_enabled         = true
  is_sparse_diskgroup_enabled     = true
  license_model                   = "LICENSE_INCLUDED"
  data_storage_size_in_tbs        = 20.0
  db_servers                      = [for db_server in data.aws_odb_db_servers_list.test.db_servers : db_server.id]
  db_node_storage_size_in_gbs     = 120.0
  memory_size_in_gbs              = 60
  data_collection_options {
    is_diagnostics_events_enabled = false
    is_health_monitoring_enabled  = false
    is_incident_logs_enabled      = false
  }
  tags = {
    "env" = "dev"
  }

}

data "aws_odb_cloud_vm_cluster" "test" {
  id = aws_odb_cloud_vm_cluster.test.id
}
`, oracleDBNetDisplayName, exaInfraDisplayName, vmcDsplayName, publicKey)
	return dsTfCodeVmCluster
}
