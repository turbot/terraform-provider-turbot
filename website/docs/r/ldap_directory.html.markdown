---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_ldap_directory"
nav:
  title: turbot_ldap_directory
---

# turbot\_ldap\_directory

The `Turbot LDAP Directory` resource adds support for ldap directories. It is used to create, manage and delete directory settings.

## Example Usage

**Creating Your First LDAP Directory**

```hcl
resource "turbot_ldap_directory" "test" {
  parent                = "tmod:@turbot/turbot#/"
  title                 = "ldap Directory"
  description           = "My first ldap directory"
  profile_id_template   = "{{profile.email}}"
  base                  = "test"
  distinguished_name    = "ldap directory"
  password              = ""
  url                   = "test.com"
}
```

## Argument Reference

The following arguments are supported:

- `base` - (Required) The base DN of the directory.
- `description` - (Optional) Brief description of the purpose and details of the directory.
- `disabled_user_filter` - (Optional) The disabled user filter of the LDAP directory to connect to
- `distinguished_name` - (Required) 
- `group_object_filter` - (Optional)
- `group_profile_id_template` - (Optional) ??
- `group_search_filter` - (Optional) The provided filter is Nunjucks rendered with `groupname` provided as a data parameter.
- `group_sync_filter` - (Optional) 
- `parent` - (Required) ID or `aka` of the parent resource.
- `profile_id_template` - (Required) A template to generate profile id for users authenticated through a ldap directory. For example, email id of the user.
- `password` - (Required) The password of the LDAP directory account
- `title` - (Required) Short descriptive name for the directory.
- `url` - (Required) The url of the LDAP directory to connect to
- `user_object_filter` - (Optional)
- `user_match_filter` - (Optional)
- `user_search_filter` - (Optional)
- `user_search_attributes` - (Optional)
- `user_canonical_name_attribute` - (Optional)
- `user_email_attribute` - (Optional)
- `user_display_name_attribute` - (Optional)
- `user_given_name_attribute` - (Optional)
- `user_family_name_attribute` - (Optional)
- `tls_enabled` - (Optional)
- `tls_server_certificate` - (Optional)
- `group_member_of_attribute` - (Optional)
- `group_membership_attribute` - (Optional)
- `connectivity_test_filter` - (Optional)
- `reject_unauthorized` - (Optional)
- `disabled_group_filter` - (Optional)
- `tags` - (Optional) Labels that can be used to manage, group, categorize, search, and save metadata for the directory.

In addition to all the arguments above, the following attributes are exported:

- `parent_akas` - A list of all `akas` for this directory's parent resource.
- `status` - Status of the ldap directory, which defaults to `Active`. Probable options are `Active`, `Inactive` and `New`.
- `id` - Unique identifier of the ldap directory.

## Import

ldap Directories can be imported using the `id`. For example,

```
terraform import turbot_ldap_directory.test 123456789012
```