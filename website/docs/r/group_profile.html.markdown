---
layout: "turbot"
title: turbot
template: Documentation
page_title: "Turbot: turbot_group_profile"
nav:
  title: turbot_group_profile
---

# turbot_group_profile

The `Turbot group profile` resource adds support for creating group profile under a directory type. It is used to create, manage and delete group profiles.

## Example Usage

**Creating Your First Group profile**

```hcl
resource "turbot_group_profile" "directory_group" {
  title             = "terraform test group profile"
  directory         = "112233445566"
  group_profile_name  = "test"
}
```

## Argument Reference

The following arguments are supported:

- `directory` - (Required) The parent directory of the group profile, either as an id, or an AKA.
- `group_profile_name` - (Required)  The unique identifier of the group profile. For new group profiles this must be unique for the parent directory.
- `title` - (Required)  The title of the group profile.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
- `status` - Status of the group profile. Valid options are `ACTIVE`, `INACTIVE`.

## Import

Turbot group profiles can be imported using the `id`. For example,

```
terraform import turbot_group_profile.test 123456789012
```