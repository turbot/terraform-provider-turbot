---
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_local_directory_user"
nav:
  title: turbot_local_directory_user
---

# turbot\_local\_directory\_user

The `Turbot Local Directory User` resource adds support for local directory users. It is used to create, manage and delete user settings.

## Example Usage

**Creating Your First Local Directory User**

```hcl
resource "turbot_local_directory_user" "test_user" {
  title               = "Local Directory User"
  email               = "xyz@turbot.com"
  display_name        = "Kai Daguerre"
  parent              = turbot_local_directory.test.id
}
```
The above example creates a local directory user called `test_user`. However, in order to create a local directory user, a local directory should be available. The following example creates a local directory and then an user for that directory.

**Creating a Local Directory and adding a User**

```hcl
resource "turbot_local_directory" "test" {
  parent                = "tmod:@turbot/turbot#/"
  title                 = "My Local"
  description           = "My first local directory"
  profile_id_template   = "{{profile.email}}"
}

resource "turbot_local_directory_user" "test_user" {
  title               = "Local Directory User"
  email               = "xyz@turbot.com"
  display_name        = "Kai Daguerre"
  parent              = turbot_local_directory.test.id
}
```

## Argument Reference

The following arguments are supported:

- `display_name` - (Required) Full display name for the new user. Usually a combination of the given name and family name.
- `email` - (Required) Email address of the new user.
- `parent` - (Required) ID or `aka` of the parent resource.
- `title` - (Required) Short descriptive name for the local directory user.
- `family_name` - (Optional) Surname of the user.
- `given_name` - (Optional) First name of the user.
- `middle_name` - (Optional) Middle name of the user.
- `picture` - (Optional) Picture of the user.
- `tags` - (Optional) Labels that can be used to manage, group, categorize, search, and save metadata for this user.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the local directory user.
- `password_timestamp` The time of the most recent change to the password field in ISO format.
- `parent_akas` -  A list of all `akas` for this user's parent resource.
- `status` -  Status of the local directory user, which defaults to `active`. Probable options are `active` and `inactive`.

## Import

Local directory user settings can be imported using the `id`. For example,

```
terraform import turbot_local_directory_user.test_user 123456789012
```
