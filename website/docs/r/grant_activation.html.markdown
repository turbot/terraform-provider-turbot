---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_grant_activation"
nav:
  title: turbot_grant_activation
---

# turbot\_grant\_activation

The `Turbot Grant Activation` resource helps activate grants for Turbot Guardrails users.

## Example Usage

**Grant Activation**

```hcl
resource "turbot_grant_activation" "test_activation" {
  resource = turbot_grant.test_grant.resource
  grant    = turbot_grant.test_grant.id
}
```

The following example creates a grant and then activates it.

```hcl
resource "turbot_grant" "test_grant" {
  resource = "tmod:@turbot/turbot#/"
  type     = "tmod:@turbot/aws#/permission/types/aws"
  level    = "tmod:@turbot/turbot-iam#/permission/levels/superuser"
  identity = turbot_profile.test.id
}

resource "turbot_grant_activation" "test_activation" {
  resource = turbot_grant.test_grant.resource
  grant    = turbot_grant.test_grant.id
}
```

## Argument Reference

The following arguments are supported:

- `resource` - (Required) The id or `aka` of the resource for which the grant is activated.
- `grant` - (Required) The `aka` of the grant to activate.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `resource_akas` - A list of all `akas` of the resource for which the grant is being activated.
- `id` - Unique identifier of the resource.

## Import

Grant Activation can be imported using the `id`. For example,

```
terraform import turbot_grant.test_grant.id 123456789012
```
