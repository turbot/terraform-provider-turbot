---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_typed_filter"
nav:
  title: turbot_typed_filter
---

# turbot\_typed\_filter

`Turbot Typed Filter` is a saved filter that a guardrail can target by reference. An `Account` filter selects the accounts a guardrail applies to, and a `Resource` filter selects the resources within them. Referencing a typed filter by its Terraform address keeps the guardrail's targeting in one place, so editing the filter changes every guardrail that uses it.

Typed filters require Guardrails 5.61 or later.

## Example Usage

**Creating an account filter and a resource filter**

```hcl
resource "turbot_typed_filter" "production_accounts" {
  title       = "Production accounts"
  type        = "Account"
  filter      = "title:production"
  description = "Accounts tagged as production"
}

resource "turbot_typed_filter" "example_bucket" {
  title  = "Example bucket"
  type   = "Resource"
  filter = "$.Name:example-bucket"
  akas   = ["example-bucket-filter"]
  tags = {
    owner = "security"
  }
}
```

## Argument Reference

The following arguments are supported:

- `title` - (Required) Short display name for the filter.
- `type` - (Required) What the filter selects: `Account` or `Resource`. Changing it destroys and recreates the filter.
- `filter` - (Required) The filter text, in the Guardrails filter syntax, for example `title:my-account` for an account or `$.Name:my-bucket` for a resource.
- `description` - (Optional) Brief description of the purpose and details of the filter.
- `akas` - (Optional) Unique identifiers of the filter. Removing all `akas` from the configuration does not remove them from the filter.
- `tags` - (Optional) Map of tags for the filter. Tags removed from the configuration are removed from the filter.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the filter.

## Import

Typed filters can be imported using the `id`. For example,

```
terraform import turbot_typed_filter.production_accounts 123456789012
```

A typed filter that a guardrail still uses cannot be deleted: destroy the guardrail, or point it at another filter, first.
