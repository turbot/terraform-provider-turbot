---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_rollout"
nav:
  title: turbot_rollout
---

# turbot_rollout

`Turbot Rollout` allows you to define and automate the rollout of one or more guardrails across one or more accounts through a series of controlled phases (e.g., preview, check, enforce, detach), with scheduled transitions and optional communications.

## Example Usage

**Automated Rollout with Scheduled Transitions**

```hcl
resource "turbot_rollout" "rollout" {

  title       = "Test Rollout Created Through Terraform"
  description = "This is a test rollout created through terraform."

  guardrails = ["356127854694770"]
  accounts   = ["350804102494366"]
  recipients = ["Account/*", "Turbot/Owner", "Turbot/Admin"]
  akas       = ["terraform_test_rollout"]
  status     = "ACTIVE"

  preview {
    start_at       = "2025-12-29T00:00:00.000Z"
    start_early_if = "no_alerts"
    start_notice   = "enabled"
  }

  check {
    start_at       = "2025-11-30T00:00:00.000Z"
    warn_at        = ["2025-10-23T00:00:00.000Z", "2025-10-29T00:00:00.000Z"]
    start_early_if = "no_alerts"
    start_notice   = "enabled"
  }

  draft {
    start_at = "2025-08-30T00:00:00.000Z"
  }
}
```

## Argument Reference

The following arguments are supported:

- `title` - (Required) A short display name for the rollout.
- `akas` - (Optional) Unique identifier of the resource.
- `description` - (Optional) Description of the rollout’s purpose.
- `guardrails` - (Optional) List of guardrail IDs or AKAs to include in the rollout.
- `accounts` - (Optional) List of account IDs or AKAs the rollout will target.
- `recipients` - (Optional) List of recipients to notify. Can include notification profiles.
- `status` - (Optional) Current status of the rollout. Must be ACTIVE to initiate phase transitions.

## Phase Blocks
At least one phase block is required. Each phase corresponds to a step in the rollout lifecycle. The valid phases are: `draft`, `preview`, `check`, `enforce`, and `detach`.

### `draft` Block

This is the initial phase and typically used for staging purposes. It supports:
- start_at – (Optional) Absolute timestamp when the rollout should transition to the draft phase.

### `preview` Block

This phase is used to introduce changes to stakeholders before enforcing them. It supports:
- start_at – (Optional) Timestamp when accounts should enter the preview phase.
- start_notice – (Optional) Whether to send welcome notices on entry. One of enabled, disabled (default).
- start_early_if – (Optional) Set to "no_alerts" to allow accounts to enter the phase early if no alerts are present.
- recipients – (Optional) Overrides the default recipients for this phase.

### `check`, `enforce`, `detach` Blocks

These phases support the full set of scheduling and communication options:
In addition to all the arguments above, the following attributes are exported:

- start_at – (Optional) Timestamp when accounts should enter the preview phase.
- warn_at – (Optional) List of timestamps to send warning notices before the phase starts.
- start_notice – (Optional) Whether to send welcome notices on entry. One of enabled, disabled (default).
- start_early_if – (Optional) Set to "no_alerts" to allow accounts to enter the phase early if no alerts are present.
- recipients – (Optional) Overrides the default recipients for this phase.

## Attributes Reference

- `parent` - The id of the rollout’s parent resource.
- `id` - Unique identifier of the resource.

## Import

Rollouts can be imported using the id. For example:

```
terraform import turbot_rollout.rollout 123456789012
```
