//Copyright © 2025, Oracle and/or its affiliates. All rights reserved.

package odb_test

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccODBDbSystemShapesListDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	dataSourceName := "data.aws_odb_db_system_shapes_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ODBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDbSystemShapesListDataSourceConfig_basic("use1-az6"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "db_system_shapes.#", "2"),
				),
			},
		},
	})
}

func testAccDbSystemShapesListDataSourceConfig_basic(availabilityZoneId string) string {
	return fmt.Sprintf(`
data "aws_odb_db_system_shapes_list" "test"{
  availability_zone_id = %[1]q
}
`, availabilityZoneId)
}
