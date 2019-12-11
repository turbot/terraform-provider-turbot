---
title: turbot_grant
template: Documentation
nav:
  title: turbot_grant
---

# turbot_grant

The `Turbot Grant` resource adds support for grants in Turbot. Turbot grant presents a cleaner and more explicit separation of duties on how users are managed in Turbot.

## Example Usage

**Creating Your First Grant**

```hcl
resource "turbot_grant" "test_grant" {
  resource = "tmod:@turbot/turbot#/"
  type     = "tmod:@turbot/aws#/permission/types/aws"
  level    = "tmod:@turbot/turbot-iam#/permission/levels/superuser"
  identity = turbot_profile.test.id
}
```
The above example creates a grant called `test_grant`. It is important to understand that a grant is for a particular profile. The following example demonstrates the process of creating a profile and granting permission.

**Creating a profile and granting permission**

```hcl
resource "turbot_profile" "test" {
  parent              = "tmod:@turbot/turbot#/"
  title               = "Snape"
  display_name        = "Severus Snape"
  email               = "severus.slytherin@hogwarts.com"
  given_name          = "Severus Snape"
  family_name         = "Snape"
  directory_pool_id   = "snapeseverus"
  status              = "Active"
  profile_id          = "170759063660234"
}

resource "turbot_grant" "test_grant" {
  resource = "tmod:@turbot/turbot#/"
  type     = "tmod:@turbot/aws#/permission/types/aws"
  level    = "tmod:@turbot/turbot-iam#/permission/levels/superuser"
  identity = turbot_profile.test.id
}
```

## Argument Reference

The following arguments are supported:

- `resource` - (Required) The id or `aka` of the resource for which permissions are being granted.
- `type` - (Required) The type of permissions being granted. This is the `aka` of a permission type resource.
- `level` - (Required) The permission level to be granted. This is the `aka` of a permission level resource.
- `identity` - (Required) The profile for which the permissions are being granted.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `resource_akas` - A list of all `akas` of the resource for which permissions are being granted.
- `permission_type_akas` - A list of all `akas` for the permission type of this grant resource.
- `permission_level_akas` - A list of all `akas` for the permission level of this grant resource.
- `identity_akas` - The `aka` of the profile for which the permissions are being granted.
- `id` - Unique identifier of the resource.

## Import

Grants can be imported using the `id`. For example,

```
terraform import turbot_grant.test_grant 123456789012
```
