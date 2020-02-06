---
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_mod"
nav:
  title: turbot_mod
---

# turbot_mod

The `Turbot Mod` resource adds support to install, update and uninstall a mod. The currently installed mod will be validated against the desired version, and the appropriate action will be taken. Removing a mod from the config will uninstall the mod.

## Example Usage

**Installing Your First Mod**

```hcl
resource "turbot_mod" "test" {
  parent    = "tmod:@turbot/turbot#/"
  org       = "turbot"
  mod       = "turbot-terraform-provider-test"
  version   = "5.0.0"
}
```

## Argument Reference

The following arguments are supported:

- `mod` - (Required) The mod to be installed, updated or uninstalled. For example, `aws-s3`.
- `org` - (Required) The parent author of the mod.
- `parent` - (Optional) Installation point for the mod in the resource hierarchy. Defaults to the Turbot root resource.
- `version` - (Optional) The version to be installed, e.g. `5.1.3`. If a semantic version range is given, e.g. `^5` then the latest available version from that range will be installed. Defaults to `*`, which is the latest available version of the mod.

**Note:** Wild cards are not accepted as inputs for pre-releases.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id` - Unique identifier of the resource.
- `version_current` - This attribute stores the version that’s currently installed (as the `version` property might be a range).
- `version_latest` - The latest version that satisfies the version requirements.
- `parent_akas` - A list of all `akas` for this mods's parent resource.
- `uri` - An unique identifier of the mod.

## Import

Mods can be imported using the `id`. For example,

```
terraform import turbot_mod.test 123456789012
```