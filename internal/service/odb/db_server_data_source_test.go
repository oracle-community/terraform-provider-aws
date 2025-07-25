//Copyright © 2025, Oracle and/or its affiliates. All rights reserved.

package odb_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	odbtypes "github.com/aws/aws-sdk-go-v2/service/odb/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/odb"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"

	tfodb "github.com/hashicorp/terraform-provider-aws/internal/service/odb"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type testDbServerDataSourceTest struct {
	exaDisplayNamePrefix string
}

var dbServerDataSourceTestEntity = testDbServerDataSourceTest{
	exaDisplayNamePrefix: "Ofake-exa",
}

// Acceptance test access AWS and cost money to run.
func TestAccODBDbServerDataSource(t *testing.T) {
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

func (testDbServerDataSourceTest) testAccCheckDbServerExists(ctx context.Context, name string, output *odb.GetDbServerOutput) resource.TestCheckFunc {
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

func (testDbServerDataSourceTest) testAccCheckDbServersDestroyed(ctx context.Context) resource.TestCheckFunc {
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

func (testDbServerDataSourceTest) findExaInfra(ctx context.Context, conn *odb.Client, id string) (*odbtypes.CloudExadataInfrastructure, error) {
	input := odb.GetCloudExadataInfrastructureInput{
		CloudExadataInfrastructureId: aws.String(id),
	}

	out, err := conn.GetCloudExadataInfrastructure(ctx, &input)
	if err != nil {
		if errs.IsA[*odbtypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: &input,
			}
		}

		return nil, err
	}

	if out == nil || out.CloudExadataInfrastructure == nil {
		return nil, tfresource.NewEmptyResultError(&input)
	}

	return out.CloudExadataInfrastructure, nil
}

func (testDbServerDataSourceTest) findDbServer(ctx context.Context, conn *odb.Client, dbServerId *string, exaInfraId *string) (*odb.GetDbServerOutput, error) {
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

func (testDbServerDataSourceTest) basic(exaInfra string) string {
	return fmt.Sprintf(`
%s

data "aws_odb_db_servers_list" "test" {
  cloud_exadata_infrastructure_id = aws_odb_cloud_exadata_infrastructure.test.id
}

data "aws_odb_db_server" "test" {
  id = data.aws_odb_db_servers_list.test.db_servers[0].id
  cloud_exadata_infrastructure_id = aws_odb_cloud_exadata_infrastructure.test.id
}
`, exaInfra)
}

func (testDbServerDataSourceTest) exaInfra(rName string) string {
	exaRes := fmt.Sprintf(`
resource "aws_odb_cloud_exadata_infrastructure" "test" {
  display_name          = "%[1]s"
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
`, rName)
	return exaRes
}

/*func (testDbServerDataSourceTest) foo(dbServerId, exaInfraId string) string {
	return fmt.Sprintf(`

data "aws_odb_db_server" "test" {
  id = %[1]q
cloud_exadata_infrastructure_id = %[2]q
}
`, dbServerId, exaInfraId)
}
*/
