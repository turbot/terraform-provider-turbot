---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_file"
nav:
  title: turbot_file
---

# turbot_file

The `Turbot file` resource allows storage of custom data within Turbot.

## Example Usage

**Creating Your First file**

```hcl
resource "turbot_file" "test" {
  parent          = "tmod:@turbot/turbot#/"
  title           = "My file"
  description     = "My first test file"
}
```

**Creating a file with content**

```hcl
resource "turbot_file" "data_file" {
  parent          = "tmod:@turbot/turbot#/"
  title           = "data file"
  description     = "This is file contains data."
  content         = <<EOF
  {
   "title": "provider_test",
   "description": "test resource"
  }EOF
}
```

## Argument Reference

The following arguments are supported:

- `content` - (Optional) Data of a file resource.
- `description` - (Optional) Brief description of the purpose and details of the file.
- `parent` - (Required) ID or `aka` of the parent resource.
- `title` - (Required) Short descriptive name for the file. This appears as the file name in the Turbot Console.
- `tags` - (Optional) Labels that can be used to manage, group, categorize, search, and save metadata for this file.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `parent_akas` - A list of all akas for this file’s parent resource.

## Import

files can be imported using the `id`. For example,

```
terraform import turbot_file.test 123456789012
```
