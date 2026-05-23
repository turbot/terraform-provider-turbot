---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_guardrail"
nav:
  title: turbot_guardrail
---

# turbot_guardrail

`Turbot Guardrail` in Turbot Governance provide a structured way to enforce compliance by grouping a control with its associated policies.

## Example Usage

**Creating Your First Guardrail**

```hcl
resource "turbot_guardrail" "aws_s3_encryption_in_transit" {
  title       = "AWS S3 S3 Bucket Encryption in Transit"
  description = "Ensure that access to Amazon S3 objects is only permitted through HTTPS, not HTTP."
  akas        = [ "aws-s3-encryption-in-transit" ]
  
  targets     = [ "tmod:@turbot/aws#/resource/types/account" ]
  controls    = [ "tmod:@turbot/aws-s3#/policy/types/encryptionInTransit" ]

  tags = {
    baseline = "required"
  }
}
```

## Argument Reference

The following arguments are supported:

- `controls` - (Required) A list of control types associated with the guardrail.
- `title` - (Required) Short display name for the guardrail.
- `akas` - (Optional) Unique identifier of the resource.
- `description` - (Optional) Brief description of the purpose and details of the guardrail.
- `tags` - (Optional) User defined label for grouping guardrails.
- `targets` - (Optional) A list of targets where the guardrail will be applied.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `color` - The color of the guardrail to create, that will be used to highlight the Guardrail in the Turbot console.
- `id` - Unique identifier of the resource.

## Import

Guardrails can be imported using the `id`. For example,

```
terraform import turbot_guardrail.test 123456789012
```
