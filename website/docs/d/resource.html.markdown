---
title: "Data Source: turbot_resource"
template: Documentation
nav:
  title: turbot_resource
---

# Data Source: turbot_resource
This data source can be used to fetch information about a specific resource.


## Example Usage

```hcl
data "turbot_resource" "test_resource" {
  id = "172720296209928"
}

output "json" {
  value = data.turbot_resource.test_resource.json_data
}
```

## Argument Reference

* `id` - (Required) The unique identifier of the resource.

## Attributes Reference

* `data` - JSON representation of the details of the resource. When parsed, it must be valid for the type schema.
* `metadata` - A set of data that describes and gives information about the data of the resource
* `akas` - A list of akas for the resource
* `tags` - User defined way of logically grouping resources.
* `json_data` - JSON representation of the full resource.
* `turbot` - JSON representation of the full data of theresource stored in Turbot.