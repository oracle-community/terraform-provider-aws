//Copyright © 2025, Oracle and/or its affiliates. All rights reserved.

package odb_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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

type listVMCListDSTest struct {
}

func TestAccListVmClusterDataSource(t *testing.T) {
	ctx := acctest.Context(t)
	var vmcListTest = listVMCListDSTest{}
	var output odb.ListCloudVmClustersOutput

	dataSourceName := "data.aws_odb_cloud_vm_clusters_list.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			vmcListTest.testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ODBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: vmcListTest.basic(),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.ComposeTestCheckFunc(func(s *terraform.State) error {
						vmcListTest.count(ctx, dataSourceName, &output)
						resource.TestCheckResourceAttr(dataSourceName, "cloud_autonomous_vm_clusters.#", strconv.Itoa(len(output.CloudVmClusters)))
						return nil
					},
					),
				),
			},
		},
	})
}

func (listVMCListDSTest) basic() string {
	config := fmt.Sprintf(`


data "aws_odb_cloud_vm_clusters_list" "test" {

}
`)
	return config
}

func (listVMCListDSTest) count(ctx context.Context, name string, list *odb.ListCloudVmClustersOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameCloudVmCluster, name, errors.New("not found"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)

		resp, err := conn.ListCloudVmClusters(ctx, &odb.ListCloudVmClustersInput{})
		if err != nil {
			return create.Error(names.ODB, create.ErrActionCheckingExistence, tfodb.ResNameCloudVmCluster, rs.Primary.ID, err)
		}

		list.CloudVmClusters = resp.CloudVmClusters

		return nil
	}
}
func (listVMCListDSTest) testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ODBClient(ctx)

	input := &odb.ListCloudVmClustersInput{}

	_, err := conn.ListCloudVmClusters(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}
