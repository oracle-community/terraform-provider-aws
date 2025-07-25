// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redshift

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkDataSource("aws_redshift_data_shares", name="Data Shares")
func newDataSharesDataSource(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSharesDataSource{}, nil
}

const (
	DSNameDataShares = "Data Shares Data Source"
)

type dataSharesDataSource struct {
	framework.DataSourceWithModel[dataSharesDataSourceModel]
}

func (d *dataSharesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_shares": framework.DataSourceComputedListOfObjectAttribute[dataSharesData](ctx),
			names.AttrID:  framework.IDAttribute(),
		},
	}
}
func (d *dataSharesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Meta().RedshiftClient(ctx)

	var data dataSharesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue(d.Meta().Region(ctx))

	paginator := redshift.NewDescribeDataSharesPaginator(conn, &redshift.DescribeDataSharesInput{})

	var out redshift.DescribeDataSharesOutput
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				create.ProblemStandardMessage(names.Redshift, create.ErrActionReading, DSNameDataShares, data.ID.String(), err),
				err.Error(),
			)
			return
		}

		if page != nil && len(page.DataShares) > 0 {
			out.DataShares = append(out.DataShares, page.DataShares...)
		}
	}

	resp.Diagnostics.Append(flex.Flatten(ctx, out, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type dataSharesDataSourceModel struct {
	framework.WithRegionModel
	DataShares fwtypes.ListNestedObjectValueOf[dataSharesData] `tfsdk:"data_shares"`
	ID         types.String                                    `tfsdk:"id"`
}

type dataSharesData struct {
	DataShareARN types.String `tfsdk:"data_share_arn"`
	ManagedBy    types.String `tfsdk:"managed_by"`
	ProducerARN  types.String `tfsdk:"producer_arn"`
}
