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

###Creating Your First LDAP Directory

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

- `base` - (Required) The BaseDN of all the requests that are sent to LDAP server. In order for your users and groups to be found in an application, they must be located underneath the base DN.
- `description` - (Optional) Brief description of the purpose and details of the directory.
- `disabled_user_filter` - (Optional) The disabled user filter of the LDAP directory to connect to
- `distinguished_name` - (Required) This is the username that Turbot will use to authenticate with the directory server after connection has been established. Mostly referred to BindDN in LDAP. 
- `group_object_filter` - (Optional) The filter string that lists all relevant groups in the LDAP server.
- `group_profile_id_template` - (Optional) A template to generate the URN of the profile for groups retrieved from this directory.
**Note**: The profile MUST be unique across all group profiles in Turbot. However, it is possible to have multiple directories map its group to the same Group-Profile.
- `group_search_filter` - (Optional) The provided filter is Nunjucks rendered with `groupname` provided as a data parameter.
- `group_sync_filter` - (Optional) Used to filter out groups of a user which Turbot should sync. If not specified, Turbot will create GroupProfiles for all groups of a user.
- `parent` - (Required) ID or `aka` of the parent resource.
- `profile_id_template` - (Required) A template to generate profile id for users authenticated through a ldap directory. For example, email id of the user.
- `password` - (Required) The password of the user specified in directory server.
- `title` - (Required) Short and descriptive name to help identify the directory server.
- `url` - (Required) The FQDN (fully qualified domain name) where the directory server is available. The FQDN also includes the protocol to be used for accessing the server.
**Note**: The host must be resolvable reachable from the network where Turbot is installed. 
- `user_object_filter` - (Optional) The filter string that lists all users in the LDAP server.
- `user_match_filter` - (Optional) This is overlaid on the `user_object_filter` to query for a specific user from the directory.
- `user_search_filter` - (Optional) This is overlaid on the `user_object_filter` to query for a sublist of all users in the directory.
- `user_search_attributes` - (Optional) This is a list of properties that will be requested from the LDAP server when requesting for a LDAP user object.
- `user_canonical_name_attribute` - (Optional) The attribute in the LDAP user object which contains the Canonical Name of the user.
- `user_email_attribute` - (Optional) The attribute in the LDAP user object which contains the user's email address.
- `user_display_name_attribute` - (Optional) The attribute in the LDAP user object which contains the user's name which will be displayed in Turbot.
- `user_given_name_attribute` - (Optional) The attribute in the LDAP user object which contains the user's Given Name.
- `user_family_name_attribute` - (Optional) The attribute in the LDAP user object which contains the user's Family Name.
- `tls_enabled` - (Optional) TLS setting for the directory server makes sure which certificates should be accepted, rejects those are not signed by a trusted CA.
- `tls_server_certificate` - (Optional) A valid and signed certificate from a trusted CA.
- `group_member_of_attribute` - (Optional) The name of the attribute which the LDAP server uses to record group memberships of an object. This is essentially the inverse of `group_membership_attribute`.
- `group_membership_attribute` - (Optional) The name of the attribute which the LDAP server uses to record membership against a group object.    
- `connectivity_test_filter` - (Optional) A filter string which will be used to test communication status with the LDAP server. 
- `reject_unauthorized` - (Optional) When TLS connection is set up, if this is set to `true`, turbot will not verify the TLS certificate that it receives from the server.
- `disabled_group_filter` - (Optional) A filter string that when queried in the context of `group_object_filter` returns disabled groups.
- `tags` - (Optional) Labels that can be used to manage, group, categorize, search, and save metadata for the directory.


In addition to all the arguments above, the following attributes are exported:

- `parent_akas` - A list of all `akas` for this directory's parent resource.
- `status` - Status of the ldap directory, which defaults to `ACTIVE`. Valid options are `ACTIVE`, `INACTIVE` and `NEW`.
- `id` - Unique identifier of the ldap directory.

## Import

ldap Directories can be imported using the `id`. For example,

```
terraform import turbot_ldap_directory.test 123456789012
```