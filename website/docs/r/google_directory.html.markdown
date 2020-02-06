---
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_google_directory"
nav:
  title: turbot_google_directory
---

# turbot\_google\_directory

The `Turbot Google Directory` resource adds support for google directory. It is used to create, manage and delete directory settings.

## Example Usage

**Creating Your First Google Directory**

```hcl
resource "turbot_google_directory" "test" {
  parent                = "tmod:@turbot/turbot#/"
  title                 = "Google @ myorg"
  profile_id_template   = "myemail@myorg.com"
  client_id             = "GoogleDirTest4"
  client_secret         = "fb-tbevaACsBKQHthzba-PH9"
  pgp_key               = "*****"
}
```

## Argument Reference

The following arguments are supported:

- `parent` - (Required) ID or `aka` of the parent resource.
- `title` - (Required) Short descriptive name for the directory.
- `profile_id_template` - (Required) A template to generate profile id for users authenticated through a google directory. For example, email id of the user.
- `client_id` - (Required) Client ID provided by Google.
- `client_secret` - (Required) Client Secret provided by Google.
- `group_id_template` - (Optional) In case of a group profile, this template generates profile id for users authenticated through this directory. For example, email id of the group.
- `login_name_template` - (Optional) A template used to render login name for the users of this directory.
- `description` - (Optional) Brief description of the purpose and details of the directory.
- `hosted_name` - (Optional) Domain name of the organization.
- `tags` - (Optional) Labels that can be used to manage, group, categorize, search, and save metadata for this directory.
- `pgp_key` - (Optional) A base-64 encoded PGP public key, applies on resource creation. If specified, the resource is encrypted in the state file with the key specified.
- `pool_id` - (Optional) Pool id associated with Google directory.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `parent_akas` - A list of all `akas` for this directory's parent resource.
- `status` -  Status of this directory, which defaults to `Active`. Probable options are `Active`, `Inactive` and `New`.
- `directory_type` - Type of the directory. For example, `google`.
- `key_fingerprint` - Unique sequence of letters and numbers used to identify a key.
- `id` - Unique identifier of the google directory.

## Import

Google Directory can be imported using the `id`. For example,

```
terraform import turbot_google_directory.test 123456789012
```