---
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_folder"
nav:
  title: turbot_folder
---

# turbot_folder

The `Turbot Folder` resource adds support for folders. Turbot folders define a top level hierarchy to arrange resources and their policies.

## Example Usage

**Creating Your First Folder**

```hcl
resource "turbot_folder" "test" {
  parent          = "tmod:@turbot/turbot#/"
  title           = "My Folder"
  description     = "My first test folder"
}
```

**Creating a Child Folder**

```hcl
resource "turbot_folder" "parent" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "Parent Folder"
  description = "This is my parent folder."
  tags = {
    "Name"        = "Provider Test"
    "Environment" = "foo"
  }
}

resource "turbot_folder" "child" {
  parent      = turbot_folder.parent.id
  title       = "Child Folder"
  description = "This is my child folder."
  tags = {
    "Name"        = "Provider Test"
    "Environment" = "foo"
  }
}
```

## Argument Reference

The following arguments are supported:

- `description` - (Required) Brief description of the purpose and details of the folder.
- `parent` - (Required) ID or `aka` of the parent resource.
- `title` - (Required) Short descriptive name for the folder. This appears as the folder name in the Turbot Console.
- `tags` - (Optional) Labels that can be used to manage, group, categorize, search, and save metadata for this folder.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `parent_akas` - A list of all akas for this folder’s parent resource.

## Import

Folders can be imported using the `id`. For example,

```
terraform import turbot_folder.test 123456789012
```
