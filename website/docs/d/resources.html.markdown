---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_resources"
nav:
  title: turbot_resources
---

# Data Source: turbot_resources
This data source can be used to search the CMDB for resources dynamically.

## Example Usage

```hcl
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test folder for guardrails terraform provider"
	tags = {
			Name = "terraform-test"
		}
}

data "turbot_resources" "test_resource" {
  filter = "tags:Name=terraform-test resourceType:folder"
}
```

## Argument Reference

* `filter` - (Required) The filter to apply to the list of resources.

## Attributes Reference

* `ids` - List of resource IDs matching the filter.