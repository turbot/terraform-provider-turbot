---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_resource"
nav:
  title: turbot_resource
---

# Data Source: turbot_resource
This data source can be used to fetch information about a specific resource.


## Example Usage

```hcl
data "turbot_resource" "test_resource" {
  id = "arn:aws:s3:::my-test"
}

output "json" {
  value = "${data.turbot_resource.test_resource}".data
}
```

## Argument Reference

* `id` - (Required) The unique identifier of the resource.

## Attributes Reference

* `data` - JSON representation of the details of the resource. When parsed, it must be valid for the type schema.
* `metadata` - A set of data that describes and gives information about the data of the resource
* `akas` - A list of akas for the resource
* `tags` - User defined way of logically grouping resources.
* `turbot` - JSON representation of turbot data of the resource.