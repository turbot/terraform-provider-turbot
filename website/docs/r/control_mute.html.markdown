---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_control_mute"
nav:
  title: turbot_control_mute
---

# turbot_control_mute

Mute controls if you want to ignore them. The `turbot_control_mute` allows muting a control to help streamline operations without compromising security policies.

## Example Usage

**Suppress alarms and errors for an `AWS > EC2 > Instance > Approved` control**

```hcl
resource "turbot_control_mute" "mute_instance_approved" {
  resource     = "arn:aws:ec2:us-east-1:123456789012:instance/i-0a2f40f8ac841fa32"
  control_type = "tmod:@turbot/aws-ec2#/control/types/instanceCmdb"
  note         = "Muting the control till alarm/error is resolved."
  to_timestamp = "2025-10-26T06:07:19.652Z"
  until_states = ["ok"]
}
```

**Suppress error for an `AWS > S3 > Bucket > Tags` control**

```hcl
resource "turbot_control_mute" "mute_bucket_tags" {
  resource     = "arn:aws:s3:::bucket-for-turbot-tags-validation"
  control_type = "tmod:@turbot/aws-s3#/control/types/bucketTags"
  note         = "Muting the control till error is resolved."
  to_timestamp = "2025-10-26T06:07:19.652Z"
  until_states = ["ok", "alarm", "skipped"]
}
```

## Argument Reference

The following arguments are supported:

- `control_id` - (Optional) The `id` or `aka` of the control to mute. **Note:** Either `control_id` or the combination of `control_type` and `resource` must be provided.
- `control_type` - (Optional) The `id` or `aka` of the control type to be muted. Must be used in combination with `resource` if `control_id` is not provided.
- `note` - (Optional) A `note` explaining the reason for muting the control.
- `resource` - (Optional) The `id` or `aka` of the resource where the control is available. Must be used in combination with `control_type` if `control_id` is not provided.
- `to_timestamp` - (Optional) A timestamp, in ISO8601 format specifying when the mute should end.
- `until_states` - (Optional) A list of control states where the mute will not apply.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the control to mute.
- `state` - The state of the specified control.

## Import

Control Mute can be imported using the `id`. For example,

```
terraform import turbot_control_mute.test 123456789012
```
