---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_policy_pack_attachment"
nav:
  title: turbot_policy_pack_attachment
---

# turbot\_policy\_pack\_attachment

The `Turbot Policy Pack Attachment` resource attaches the policy pack to specific Turbot Guardrails resources.

## Example Usage

**Creating Your Policy Pack**

```hcl
resource "turbot_policy_pack" "policy_pack" {
  title = "My policy pack"
}
```

**Creating Your Resource**

```hcl
resource "turbot_resource" "my_resource" {
  parent = "tmod:@turbot/turbot#/"
  type   = "tmod:@turbot/aws#/resource/types/account"
  data   = <<EOT
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

**Attaching Your Policy Pack to the Resource**

```hcl
resource "turbot_policy_pack_attachment" "test" {
  resource    = "${turbot_resource.my_resource.id}"
  policy_pack = "${turbot_policy_pack.policy_pack.id}"
}
```
The above example attaches a policy pack (`policy_pack`) to a resource (`my_resource`). It is important to understand that you have a policy pack and a resource for the attachment process. The following example provides a systematic approach of creating a policy pack attachment.

## Argument Reference

The following arguments are supported:

- `policy_pack` - (Required) The id of the policy pack to be attached.
- `resource` - (Required) The id of the resource to which a policy pack will be attached.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
