---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "AWS: aws_ec2_transit_gateway_route_table_associations"
description: |-
   Provides information for multiple EC2 Transit Gateway Route Table Associations
---

# Data Source: aws_ec2_transit_gateway_route_table_associations

Provides information for multiple EC2 Transit Gateway Route Table Associations, such as their identifiers.

## Example Usage

### By Transit Gateway Identifier

```terraform
data "aws_ec2_transit_gateway_route_table_associations" "example" {
  transit_gateway_route_table_id = aws_ec2_transit_gateway_route_table.example.id
}
```

## Argument Reference

The following arguments are required:

* `transit_gateway_route_table_id` - (Required) Identifier of EC2 Transit Gateway Route Table.

The following arguments are optional:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `filter` - (Optional) Custom filter block as described below.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) Name of the field to filter by, as defined by
  [the underlying AWS API](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_GetTransitGatewayRouteTableAssociations.html).

* `values` - (Required) Set of values that are accepted for the given field.
  A Transit Gateway Route Table will be selected if any one of the given values matches.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - AWS Region.
* `ids` - Set of Transit Gateway Route Table Association identifiers.
