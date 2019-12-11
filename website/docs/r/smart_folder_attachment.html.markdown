---
title: turbot_smart_folder_attachment
template: Documentation
nav:
  title: turbot_smart_folder_attachment
---

# turbot\_smart\_folder\_attachment

The `Turbot Smart Folder Attachment` resource attaches the smart folder to specific Turbot resources.



## Example Usage

**Creating a Smart Folder Attachment**

```hcl
resource "turbot_smart_folder_attachment" "test" {
  resource     = "167225763707951"
  smart_folder = "171222424857954"
}
```
The above example attaches a smart folder (`id: 171222424857954`) to a resource (`id: 167225763707951`). It is important to understand that you have a smart folder and a resource for the attachment process. The following example provides a systematic approach of creating a smart folder attachment.

**Creating Your Smart Folder**

```hcl
resource "turbot_smart_folder" "test" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "My smart folder"
}
```

**Creating Your Resource**

```hcl
resource "turbot_resource" "my_resource" {
  parent      = "tmod:@turbot/turbot#/"
  type        = "tmod:@turbot/aws#/resource/types/account"

  payload =  <<EOF
  {
    "Id": "123456789000",
    "turbot": {}
  }
  EOF
}
```

**Attaching Your Smart Folder to the Resource**

```hcl
resource "turbot_smart_folder_attachment" "test" {
  resource     = "167225763707951"
  smart_folder = "171222424857954"
}
```

## Argument Reference

The following arguments are supported:

- `resource` - (Required) The id of the resource to which a smart folder will be attached.
- `smart_folder` - (Required) The id of the smart folder to be attached.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
