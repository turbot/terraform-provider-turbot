---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_resource"
nav:
  title: turbot_resource
---

# turbot_resource

The `turbot_resource` defines a resource in Turbot. Typically it is used to define the top level for a set of discoverable resources (e.g. an AWS account).

## Example Usage

**Creating Your First Resource**

```hcl
resource "turbot_resource" "my_resource" {
  parent   = "tmod:@turbot/turbot#/"
  type     = "tmod:@turbot/aws#/resource/types/account"
  akas     = ["arn:aws:::123456789012"]
  data     = <<EOT
{
  "Id": "123456789012",
}
EOT
  metadata = <<EOT
{
  "aws": {
    "accountId": "123456789012",
    "partition": "aws"
  }
}
EOT
}
```

## Argument Reference

The following arguments are supported:

- `parent` - (Required) The `id` or `aka` of the level at which the Turbot resource will be created.
- `type` - (Required) Defines the type of the resource to be created.
- `data` - (Required) JSON representation of the details of the resource. When parsed, it must be valid for the `type` schema.
- `metadata` - (Optional) A set of data that describes and gives information about the data of the resource.
- `akas` - (Optional) Unique identifier of the resource.
- `tags` - (Optional) User defined label for grouping resources.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
- `parent_akas` - A list of all `akas` for the Turbot resource's parent resource.

## Import

Resources can be imported using the `id`. For example,

```
terraform import turbot_resource.my_account 123456789012
```
