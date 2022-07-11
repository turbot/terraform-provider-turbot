---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_watch"
nav:
  title: turbot_watch
---

# turbot_watch

`Turbot Watch` helps in monitoring of specific events, such as control updates, grant activations, and resource creation.

## Example Usage

**Creating Your First Watch**

```hcl
resource "turbot_watch" "watch" {
  resource = "185423120545381"
  action   = "tmod:@turbot/firehose-aws-sns#/action/types/router"
  filters  = [
    "level:self,descendant notificationType:active_grants_deleted"
  ]
}
```

## Argument Reference

The following arguments are supported:

- `resource` - (Required) The resource to create the Watch for, either a Turbot ID or AKA.
- `filters` - (Required) A valid reverse filter to determine which notifications to process.
- `action` - (Required) The action the Watch takes when it finds a match, either a Turbot ID or URI.
- `favorite` - (Optional) Favorite to associate the Watch with, as a Turbot ID.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `description` - A brief description for the watch.
- `handler` - The handler object for the watch.
- `id` - Unique identifier of the watch.

## Import

Watches can be imported using the `id`. For example,

```
terraform import turbot_watch.test 123456789012
```
