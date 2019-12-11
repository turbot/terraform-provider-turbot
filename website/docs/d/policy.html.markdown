---
title: "Data Source: turbot_policy_value"
template: Documentation
nav:
  title: turbot_policy_value
---

# Data Source: turbot\_policy\_value

This data source can be used to fetch information about a specific policy
setting.

## Example Usage

A simple example to extract the value of a policy.

```hcl
data "turbot_policy_value" "example" {
  type      = "tmod:@turbot/turbot-iam#/policy/types/permissions"
  resource  = "172720296209928"
}
```
Here is another example wherein the value of a turbot policy is used to set another policy on a folder.

```hcl
data "turbot_policy_value" "example" {
  type            = "tmod:@turbot/turbot-iam#/policy/types/permissions"
  resource        = "172720296209928"
}

output "op1" {
  value = "data.turbot_policy_value.example.value"
}

resource "turbot_folder" "parent" {
  parent        = "tmod:@turbot/turbot#/"
  title         = "Data Source"
  description   = "Testing the policy data source of Turbot"
}
resource "turbot_policy_setting" "test_policy" {
  resource      = "turbot_folder.parent.id"
  type          = "tmod:@turbot/turbot-iam#/policy/types/permissions"
  value         = "data.turbot_policy_value.example.value"
  precedence    = "must"
}
```

## Argument Reference

* `type` - (Required) The unique identifier of the policy for which the value needs to be extracted.
* `resource` - (Required) The unique ID of the resource at the level of which the information needs to be fetched.


## Attributes Reference

* `value` - The value that the policy is set to.
* `value_source` - The values for the policy derived from the template.
* `precedence` - The priority level of the policy.
* `state` - The final state of the set policy.
* `reason` - Message explaining the state of the set policy.
* `details` - Additional information regarding the set policy.
* `setting_id` - The unique id of the the policy setting.