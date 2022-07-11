---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_smart_folder"
nav:
  title: turbot_smart_folder
---

# turbot\_smart\_folder

`Turbot Smart Folder` allows resources from across the hierarchy to be organized together as a group.

## Example Usage

**Creating Your First Smart Folder**

```hcl
resource "turbot_smart_folder" "folder" {
  parent  = "tmod:@turbot/turbot#/"
  title   = "My smart folder"
}
```

## Argument Reference

The following arguments are supported:

- `parent` - (Required) The `id` or `aka` of the level at which the smart folder will be created.
- `title` - (Required) Short display name for the smart folder.
- `description` - (Optional) Brief description of the purpose and details of the smart folder.
- `filter` - (Optional) A query syntax to identify the resources onto which the smart folder will automatically get attached.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `parent_akas` - A list of all `akas` for this smart folder’s parent resource.
- `id` - Unique identifier of the resource.

## Import

Smart Folders can be imported using the `id`. For example,

```
terraform import turbot_smart_folder.test 123456789012
```
