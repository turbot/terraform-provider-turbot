---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_guardrail"
nav:
  title: turbot_guardrail
---

# turbot\_guardrail

`Turbot Guardrail` is an instance of a guardrail type: a mode, an optional action, the type's settings, and the accounts and resources it applies to.

Guardrails require Guardrails 5.61 or later.

## Example Usage

**A guardrail with fixed settings**

```hcl
resource "turbot_guardrail" "bucket_encryption" {
  type            = "tmod:@turbot/aws-s3#/guardrail/types/bucketEncryptionAtRestGuardrail"
  title           = "Bucket encryption at rest"
  description     = "S3 buckets in production accounts use S3 managed keys"
  mode            = "Check"
  grace_period    = "14d"
  account_filter  = "tags:environment=production"
  resource_filter = "$.Name:/^prod-/"

  action {
    value = "tmod:@turbot/aws-s3#/action/types/bucketUpdateEncryptionAtRest"
  }

  setting {
    type  = "encryptionLevel"
    value = "S3 managed key"
  }
}
```

The type's other settings, `acceptStrongerEncryption` and `customerManagedKey`, have no `setting` block, so they take the type's defaults.

**A guardrail with a calculated setting, targeting typed filters**

```hcl
resource "turbot_typed_filter" "production_accounts" {
  title  = "Production accounts"
  type   = "Account"
  filter = "tags:environment=production"
}

resource "turbot_guardrail" "bucket_versioning" {
  type              = "tmod:@turbot/aws-s3#/guardrail/types/bucketVersioningGuardrail"
  title             = "Bucket versioning"
  mode              = "Check"
  account_filter_id = turbot_typed_filter.production_accounts.id

  setting {
    type           = "versioning"
    template_input = <<-EOT
      {
        resource {
          tags
        }
      }
    EOT
    template       = <<-EOT
      {%- if $.resource.tags.archive == "true" -%}
      Disabled
      {%- else -%}
      Enabled
      {%- endif -%}
    EOT
  }
}
```

## Argument Reference

The following arguments are supported:

- `type` - (Required) The URI of the guardrail type, for example `tmod:@turbot/aws-s3#/guardrail/types/bucketVersioningGuardrail`. Changing it creates a new guardrail.
- `title` - (Required) Short display name for the guardrail.
- `mode` - (Required) The enforcement mode, one of the type's modes: `Skip`, `Draft`, `Check`, `Enforce if new` or `Enforce all`. A static value; the API rejects a calculated mode.
- `description` - (Optional) Brief description of the guardrail.
- `akas` - (Optional) Unique identifiers of the guardrail. Removing all `akas` from the configuration does not remove them from the guardrail.
- `tags` - (Optional) Tags for the guardrail.
- `grace_period` - (Optional) How long to wait before enforcing, for example `14d`.
- `action` - (Optional) The action the guardrail takes, as one block holding either `value`, the full URI of one of the type's action types, or `template` and `template_input` for an action calculated per resource. Once a guardrail has an action the API cannot remove it, so removing this block is an error at plan time.
- `account_filter` - (Optional) A Turbot filter selecting the accounts the guardrail applies to. Conflicts with `account_filter_id`.
- `account_filter_id` - (Optional) The `id` of an `Account` [turbot_typed_filter](typed_filter.html) selecting the accounts. Conflicts with `account_filter`. Account targeting can be switched between `account_filter` and `account_filter_id`, but the API cannot clear it, so removing both is an error at plan time.
- `resource_filter` - (Optional) A Turbot filter selecting which resources in the targeted accounts the guardrail governs. Conflicts with `resource_filter_id`. With neither set, the guardrail governs every resource of its type.
- `resource_filter_id` - (Optional) The `id` of a `Resource` [turbot_typed_filter](typed_filter.html). Conflicts with `resource_filter`.
- `setting` - (Optional) A value for one of the type's settings. Repeat the block for each setting you set; each setting may appear once.

### setting

- `type` - (Required) The setting's short name, for example `versioning` or `encryptionLevel`. An unknown name is an error that lists the type's settings.
- `value` - (Optional) A fixed value, as a YAML or JSON string, for example `"Enabled"` or `"true"`.
- `template` - (Optional) A Nunjucks template that calculates the value per resource. Set either `value`, or `template` and `template_input`.
- `template_input` - (Optional) The GraphQL query whose result the template reads.

A setting with no `setting` block takes the type's default: creating the guardrail sends the default, and removing a block later writes the default back. A setting whose type defines no default is left unset on create, and its block cannot be removed once set: give it a value instead.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the guardrail.

## Import

Guardrails can be imported using the `id`. For example,

```
terraform import turbot_guardrail.test 123456789012
```

A setting whose value is the type's default is not imported, which matches configuration that has no block for it.
