---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_local_directory"
nav:
  title: turbot_local_directory
---

# turbot\_local\_directory

The `Turbot Local Directory` resource adds support for local directories. It is used to create and delete directory settings.

## Example Usage

**Creating Your First Local Directory**

```hcl
resource "turbot_local_directory" "test" {
  parent                = "tmod:@turbot/turbot#/"
  title                 = "Local Directory"
  description           = "My first local directory"
  profile_id_template   = "{{profile.email}}"
}
```

## Argument Reference

The following arguments are supported:

- `parent` - (Required) ID or `aka` of the parent resource.
- `profile_id_template` - (Required) A template to generate profile id for users authenticated through a local directory. For example, email id of the user.
- `title` - (Required) Short descriptive name for the directory.
- `description` - (Optional) Brief description of the purpose and details of the directory.
- `tags` - (Optional) Labels that can be used to manage, group, categorize, search, and save metadata for the directory.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `parent_akas` - A list of all `akas` for this directory's parent resource.
- `status` - Status of the local directory, which defaults to `Active`. Probable options are `Active`, `Inactive` and `New`.
- `directory_type` - Type of the directory. For example, `local`.
- `id` - Unique identifier of the local directory.

## Import

Local Directories can be imported using the `id`. For example,

```
terraform import turbot_local_directory.test 123456789012
```