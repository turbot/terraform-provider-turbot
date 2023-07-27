---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_smart_folder_attachment"
nav:
  title: turbot_smart_folder_attachment
---

# turbot\_smart\_folder\_attachment

The `Turbot Smart Folder Attachment` resource attaches the smart folder to specific Turbot Guardrails resources.



## Example Usage

**Creating Your Smart Folder**

```hcl
resource "turbot_smart_folder" "smart_folder" {
  parent  = "tmod:@turbot/turbot#/"
  title   = "My smart folder"
}
```

**Creating Your Resource**

```hcl
resource "turbot_resource" "my_resource" {
  parent   = "tmod:@turbot/turbot#/"
  type     = "tmod:@turbot/aws#/resource/types/account"
  data     = <<EOT
{
  "Id": "123456789012",
}
EOT
  metadata = <<EOT
{
  "aws": {
    "accountId": "123456789012",
    "partition": "aws"
  }
}
EOT
}
```

**Attaching Your Smart Folder to the Resource**

```hcl
resource "turbot_smart_folder_attachment" "test" {
  resource     = "${turbot_resource.my_resource.id}"
  smart_folder = "${turbot_smart_folder.smart_folder.id}"
}
```
The above example attaches a smart folder (`smart_folder`) to a resource (`my_resource`). It is important to understand that you have a smart folder and a resource for the attachment process. The following example provides a systematic approach of creating a smart folder attachment.

## Argument Reference

The following arguments are supported:

- `resource` - (Required) The id of the resource to which a smart folder will be attached.
- `smart_folder` - (Required) The id of the smart folder to be attached.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
