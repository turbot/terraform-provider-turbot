---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_control"
nav:
  title: turbot_control
---

# Data Source: turbot\_control

This data source can be used to fetch information about a specific control.

## Example Usage

A simple example to extract the status of a control.

```hcl
data "turbot_control" "test" {
  id      = "tmod:@turbot/aws-s3#/control/types/bucketVersioning"
}

output "json" {
  value = "${data.turbot_control.test}".status
}
```
Here is another example wherein, we can fetch control data using uri and resource.

```hcl
data "turbot_control_value" "example" {
  uri      = "tmod:@turbot/aws-ec2#/control/types/instanceDiscovery"
  resource  = 'arn:aws::ap-northeast-1:112233445566'
}

output "json" {
  value = "${data.turbot_control_value.example}".status
}
```

## Argument Reference

* `id` - (Optional) The unique identifier of the control for which the value needs to be extracted.
* `uri` - (Optional) The unique identifier of the control for which the value needs to be extracted.
* `resource` - (Optional) The unique ID of the resource at the level of which the information needs to be fetched.

## Attributes Reference

* `state` - The final state of the set control.
* `reason` - Message explaining the state of the set control.
* `details` - Additional information regarding the set control.
* `tags` - User defined way of logically grouping resources.
* `turbot` - JSON representation of turbot data of the resource.