//Copyright © 2025, Oracle and/or its affiliates. All rights reserved.

package odb_test

import (
	"context"
	"errors"
	"fmt"
	odbtypes "github.com/aws/aws-sdk-go-v2/service/odb/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/odb"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/names"

	tfodb "github.com/hashicorp/terraform-provider-aws/internal/service/odb"
)

type odbNwkPeeringResourceTest struct {
	vpcNamePrefix               string
	odbPeeringDisplayNamePrefix string
	odbNwkDisplayNamePrefix     string
}

var odbPeeringTestResource = odbNwkPeeringResourceTest{
	vpcNamePrefix:               "vpc",
	odbPeeringDisplayNamePrefix: "odb-peering",
	odbNwkDisplayNamePrefix:     "odb-net",
}

func TestAccODBNetworkPeeringConnection_basic(t *testing.T) {
	ctx := acctest.Context(t)

	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var odbPeeringResource odb.GetOdbPeeringConnectionOutput
	odbPeeringDisplayName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.odbPeeringDisplayNamePrefix)
	vpcName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.vpcNamePrefix)
	odbNetName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.odbNwkDisplayNamePrefix)
	resourceName := "aws_odb_network_peering_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			//acctest.PreCheckPartitionHasService(t, names.ODBEndpointID)
			odbPeeringTestResource.testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ODBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             odbPeeringTestResource.testAccCheckNetworkPeeringConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: odbPeeringTestResource.basicConfig(vpcName, odbNetName, odbPeeringDisplayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkPeeringConnectionExists(ctx, resourceName, &odbPeeringResource),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "dev"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccODBNetworkPeeringConnectionAddRemoveTag(t *testing.T) {
	ctx := acctest.Context(t)

	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var odbPeeringResource odb.GetOdbPeeringConnectionOutput
	odbPeeringDisplayName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.odbPeeringDisplayNamePrefix)
	//vpcName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.vpcNamePrefix)
	odbNetName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.odbNwkDisplayNamePrefix)
	resourceName := "aws_odb_network_peering_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			//acctest.PreCheckPartitionHasService(t, names.ODBEndpointID)
			odbPeeringTestResource.testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ODBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             odbPeeringTestResource.testAccCheckNetworkPeeringConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: odbPeeringTestResource.basicConfigWithVPC("vpc-084bc7dd335e156cc", odbNetName, odbPeeringDisplayName),
				//odbPeeringTestResource.basicConfig(vpcName, odbNetName, odbPeeringDisplayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkPeeringConnectionExists(ctx, resourceName, &odbPeeringResource),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "dev"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: odbPeeringTestResource.basicConfigWithVPCWithNoTag("vpc-084bc7dd335e156cc", odbNetName, odbPeeringDisplayName),
				//odbPeeringTestResource.basicConfig(vpcName, odbNetName, odbPeeringDisplayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkPeeringConnectionExists(ctx, resourceName, &odbPeeringResource),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccODBNetworkPeeringConnection_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var odbPeering odb.GetOdbPeeringConnectionOutput
	odbPeeringDisplayName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.odbPeeringDisplayNamePrefix)
	vpcName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.vpcNamePrefix)
	odbNetDisplayName := sdkacctest.RandomWithPrefix(odbPeeringTestResource.odbPeeringDisplayNamePrefix)
	resourceName := "aws_odb_network_peering_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			//acctest.PreCheckPartitionHasService(t, names.ODBEndpointID)
			odbPeeringTestResource.testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ODBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		//CheckDestroy:             odbPeeringTestResource.testAccCheckNetworkPeeringConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: odbPeeringTestResource.basicConfig(vpcName, odbNetDisplayName, odbPeeringDisplayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkPeeringConnectionExists(ctx, resourceName, &odbPeering),

					//acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfodb.odbPeering, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func (odbNwkPeeringResourceTest) testAccCheckNetworkPeeringConnectionDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_odb_network_peering_connection" {
				continue
			}

			_, err := odbPeeringTestResource.findOdbPeering(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				return nil
			}
			if err != nil {
				return create.Error(names.ODB, create.ErrActionCheckingDestroyed, tfodb.ResNameNetworkPeeringConnection, rs.Primary.ID, err)
			}

			return create.Error(names.ODB, create.ErrActionCheckingDestroyed, tfodb.ResNameNetworkPeeringConnection, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckNetworkPeeringConnectionExists(ctx context.Context, name string, odbPeeringConnection *odb.GetOdbPeeringConnectionOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameNetworkPeeringConnection, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameNetworkPeeringConnection, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)

		resp, err := odbPeeringTestResource.findOdbPeering(ctx, conn, rs.Primary.ID)
		if err != nil {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameNetworkPeeringConnection, rs.Primary.ID, err)
		}

		*odbPeeringConnection = *resp

		return nil
	}
}

func (odbNwkPeeringResourceTest) testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)

	input := &odb.ListOdbPeeringConnectionsInput{}

	_, err := conn.ListOdbPeeringConnections(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

/*func testAccCheckNetworkPeeringConnectionNotRecreated(before, after *odb.GetOdbPeeringConnectionOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.ToString(before.NetworkPeeringConnectionId), aws.ToString(after.NetworkPeeringConnectionId); before != after {
			return create.Error(names.ODB, create.ErrActionCheckingNotRecreated, tfodb.ResNameNetworkPeeringConnection, aws.ToString(before.NetworkPeeringConnectionId), errors.New("recreated"))
		}

		return nil
	}
}*/

func (odbNwkPeeringResourceTest) findOdbPeering(ctx context.Context, conn *odb.Client, id string) (output *odb.GetOdbPeeringConnectionOutput, err error) {
	input := odb.GetOdbPeeringConnectionInput{
		OdbPeeringConnectionId: &id,
	}
	out, err := conn.GetOdbPeeringConnection(ctx, &input)
	if err != nil {
		if errs.IsA[*odbtypes.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: &input,
			}
		}
		return nil, err
	}
	if out == nil {
		return nil, errors.New("odb Network Peering Connection resource can not be nil")
	}
	return out, nil
}

func (odbNwkPeeringResourceTest) basicConfig(vpcName, odbNetName, odbPeeringName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"

  tags = {
    Name = %[1]q
  }
}

resource "aws_odb_network" "test" {
  display_name          = %[2]q
  availability_zone_id = "use1-az6"
  client_subnet_cidr   = "10.2.0.0/24"
  backup_subnet_cidr   = "10.2.1.0/24"
  s3_access = "DISABLED"
  zero_etl_access = "DISABLED"
}

resource "aws_odb_network_peering_connection" "test" {
  display_name = %[3]q
  odb_network_id = aws_odb_network.test.id
  peer_network_id = aws_vpc.test.id
  tags = {
    "env"="dev"
  }
}
`, vpcName, odbNetName, odbPeeringName)
}

func (odbNwkPeeringResourceTest) basicConfigWithVPC(vpcName, odbNetName, odbPeeringName string) string {
	return fmt.Sprintf(`


resource "aws_odb_network" "test" {
  display_name          = %[2]q
  availability_zone_id = "use1-az6"
  client_subnet_cidr   = "10.2.0.0/24"
  backup_subnet_cidr   = "10.2.1.0/24"
  s3_access = "DISABLED"
  zero_etl_access = "DISABLED"
}

resource "aws_odb_network_peering_connection" "test" {
  display_name = %[3]q
  odb_network_id = aws_odb_network.test.id
  peer_network_id = %[1]q
  tags = {
    "env"="dev"
  }
}
`, vpcName, odbNetName, odbPeeringName)
}

func (odbNwkPeeringResourceTest) basicConfigWithVPCWithNoTag(vpcName, odbNetName, odbPeeringName string) string {
	return fmt.Sprintf(`


resource "aws_odb_network" "test" {
  display_name          = %[2]q
  availability_zone_id = "use1-az6"
  client_subnet_cidr   = "10.2.0.0/24"
  backup_subnet_cidr   = "10.2.1.0/24"
  s3_access = "DISABLED"
  zero_etl_access = "DISABLED"
}

resource "aws_odb_network_peering_connection" "test" {
  display_name = %[3]q
  odb_network_id = aws_odb_network.test.id
  peer_network_id = %[1]q
}
`, vpcName, odbNetName, odbPeeringName)
}
