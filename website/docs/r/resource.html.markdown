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

**Using full_data**

```hcl
resource "turbot_resource" "my_resource" {
  parent   = "tmod:@turbot/turbot#/"
  type     = "tmod:@turbot/aws#/resource/types/account"
  full_data     = <<EOT
{
  "type": "aws",
  "Id": "123456789012",
  "title": "turbot account resource"
}
EOT
}
```

**Using full_metadata**

```hcl
resource "turbot_resource" "my_resource" {
  parent   = "tmod:@turbot/turbot#/"
  type     = "tmod:@turbot/aws#/resource/types/account"
  full_metadata    = <<EOT
{
  "resource_version": "1.0.0"
  "replication": "true"
}
EOT
}
```

## Argument Reference

The following arguments are supported:

- `parent` - (Required) The identifier of the parent resource under which this resource will be created.
- `type` - (Required) Defines the type of the resource to be created.
- `data` - (Optional) JSON representation of resource properties to be managed by Terraform. The data must be valid for the resource type schema. NOTE: If additional properties are set on the resource by other means, they are ignored by Terraform.
- `metadata` - (Optional) JSON representation of resource metadata properties to be managed by Terraform. NOTE: If additional metadata properties are set on the resource by other means, they are ignored by Terraform.
- `full_data` - (Optional) JSON representation of all resource properties to be set on the resource. The data must be valid for the resource type schema. NOTE: If additional properties are set on the resource by other means, they are removed.
- `full_metadata` - (Optional) JSON representation of all resource metadata properties to be set on the resource. NOTE: If additional metadata properties are set on the resource by other means, they are removed.
- `akas` - (Optional) Unique identifier of the resource.
- `tags` - (Optional) User defined label for grouping resources.
 
**NOTE**: Only one of the `data` and `full_data` must be specified. Likewise, only one of `metadata` and `full_metadata` must be set.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
- `parent_akas` - A list of all `akas` for the Turbot resource's parent resource.

## Import

Resources can be imported using the `id`. For example,

```
terraform import turbot_resource.my_account 123456789012
```
