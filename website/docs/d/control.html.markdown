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
Here is another example wherein, we can fetch control data using the control type uri and resource id.

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

* `id` - (Optional) The id of the control.
* `uri` - (Optional) The type of the control.
* `resource` - (Optional) The unique identifier of the resource which the control is targeting.

**Note:** You must specify either the control id or the control type uri AND the resource.
## Attributes Reference

* `state` - The state of the control.
* `reason` - Message explaining the state of the control.
* `details` - Additional information regarding the control state.
* `tags` - Tags set on the control.