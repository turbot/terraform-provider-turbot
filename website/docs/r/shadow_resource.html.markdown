---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_shadow_resource"
nav:
  title: turbot_shadow_resource
---

# turbot\_shadow\_resource

Shadow resources allow management of the Turbot Guardrails representation (a shadow) of
the actual resource. A shadow resource provides two important capabilities:
  * The shadow will wait until Turbot Guardrails has discovered the actual resource before being created.
  * The shadow will be deleted before the actual resource, allowing a chance for Turbot Guardrails managed changes to be cleaned up.


## Example Usage

### Set policy directly on a new cloud resource

Setting a policy directly on a newly created resource will usually fail with
a resource not found error. The resource is created, but Turbot Guardrails has not had
a chance to discover it yet. (Seconds matter!)

A shadow resource provides a way to wait for the new resource to be discovered
by Turbot Guardrails, and then perform actions (e.g. setting a policy).

```hcl
resource "aws_s3_bucket" "my_bucket" {
  bucket   = "shadow-resource-test"
}

resource "turbot_shadow_resource" "my_bucket_shadow" {
  resource = aws_s3_bucket.my_bucket.arn
}

resource "turbot_policy_setting" "s3_bucket_versioning" {
  resource = turbot_shadow_resource.my_bucket_shadow.id
  type     = "tmod:@turbot/aws-s3#/policy/types/bucketVersioning"
  value    = "Enforce: Enabled"
}
```

### Create a shadow resource from a filter search

Sometimes the ID of a resource is not yet known, but it can be uniquely
identified via a filter search. Provided exactly 1 result is returned
from the filter, a shadow resource can be created:

```hcl
resource "turbot_shadow_resource" "my_ec2_instance_shadow" {
  filter = "$.PrivateIpAddress:10.0.0.121 resourceType:instance"
}
```


## Argument Reference

At least one of `resource` or `filter` must be specified:

- `resource` - (Optional) ID of the resource that the shadow resource will represent.
- `filter` - (Optional) Filter query matching a single resource.


## Timeouts

`turbot_shadow_resource` provides the following [Timeouts](/docs/configuration/resources.html#timeouts)
configuration options:

- `create` - (Default `5m`) How long to wait for a resource to be created.
