---
layout: "turbot"
title: turbot
template: Documentation
page_title: "Turbot: turbot_profile"
nav:
  title: turbot_profile
---

# turbot_profile

The `Turbot Profile` resource adds support for creating user profiles. It is used to create, manage and delete profile settings.

## Example Usage

###Creating Your First Profile

```hcl
resource "turbot_profile" "admin" {
  parent              = "tmod:@turbot/turbot#/"
  title               = "Admin"
  display_name        = "Hoax"
  email               = "hoax@turbot.com"
  given_name          = "Hoax benjamin"
  family_name         = "Benjamin"
  profile_id          = "170759063660234"
}
```

## Argument Reference

The following arguments are supported:

- `display_name` - (Required) The display name of the profile.
- `email` - (Required) Email ID associated with the profile.
- `family_name` - (Required) Last name of the user associated with the profile.
- `given_name` - (Required) First name of the user associated with the profile.
- `parent` - (Required) The `aka` or `id` of the level at which the profile is created.
- `profile_id` - (Required) An unique identifier of the profile.
- `title` - (Required) Name of the profile.
- `directory_pool_id` - (Optional) Pool ID for the directory in the current resource. Allows grouping of related directories e.g. SAML for authentication and LDAP for AD searching.
- `external_id` - (Optional) A link between the local directory and the profile.
- `last_login_timestamp` - (Optional) The most recent login through the profile.
- `middle_name` - (Optional) Middle name of the user associated with the profile.
- `picture` - (Optional) A valid URL which contains a picture which will be associated to the profile.
- `status` - (Optional) Status of the profile, which defaults to `Active`. Valid options are `Active` and `Inactive`.

**Note:** In case of a local directory, both the `profile_id` and `external_id` are required parameters.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
- `parent_akas` - A list of all `akas` for this Turbot Guardrails profiles's parent resource.

## Import

Turbot Guardrails profiles can be imported using the `id`. For example,

```
terraform import turbot_folder.admin 123456789012
```