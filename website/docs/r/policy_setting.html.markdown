---
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_policy_setting"
nav:
  title: turbot_policy_setting
---

# turbot\_policy\_setting

The `Turbot Policy Setting` resource adds support for setting policies for resources. It is used to create, update and delete policy settings.

## Example Usage

**Creating Your First Folder**

```hcl
resource "turbot_folder" "test" {
  parent          = "tmod:@turbot/turbot#/"
  title           = "My Folder"
  description     = "My first test folder"
}
```

**Setting Your First Policy**

```hcl
resource "turbot_policy_setting" "test_policy" {
  resource    = turbot_folder.parent.id
  type        = "tmod:@turbot/turbot-iam#/policy/types/permissions"
}
```

**Setting Your Policy Using Nunjucks Template**

```hcl
resource "turbot_policy_setting" "template_policy" {
  resource          = "arn:aws::eu-west-2:650022101893"
  type              = "tmod:@turbot/aws#/policy/types/accountStack"
  template_input    = "{ account{ Id } }"
  template          = "{% if $.account.Id == '650022101893' %}Skip{% else %}'Check: Configured'{% endif %}"
  precedence        = "must"
}
```

## Argument Reference

The following arguments are supported:

- `type` - (Required) The `aka` of the policy type to be created. This is represented by `uri` which can be found out from the overview section of the desired policy.
- `resource` - (Required) The `aka` of the resource.
- `note` - (Optional) Additional notes, if desired.
- `precedence` - (Optional) Determines whether the policy setting should be `required` or `recommended`. Defaults to `required`.
- `template` - (Optional) Nunjucks template that is used to render the policy.
- `template_input` - (Optional) A GraphQL query required as the input for the `template`.
- `valid_from_timestamp` - (Optional) The start of a specific time period for which the policy setting is valid.
- `valid_to_timestamp` - (Optional) The expiration date of a policy value.
- `value` - (Optional) Value of the policy. This could either be the value of the setting or a `yaml` string representing the setting.
- `pgp_key` - (Optional) A base-64 encoded PGP public key, applies on resource creation. If specified, the resource is encrypted in the state file with the key specified.


## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
- `value_source` - The YAML representation of the policy.
- `value_key_fingerprint` -  Value of the fingerprint used to identify a key
- `value_source_key_fingerprint` - The source of the value of the key fingerprint.
- `value_source_used` - The YAML representation of the policy that is in use.

## Import

Policy settings can be imported using the `id`. For example,

```
terraform import turbot_policy_setting.s3_encryption_at_rest 123456789012
```
