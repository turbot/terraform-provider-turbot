---
layout: "turbot"
title: "turbot"
template: Documentation
page_title: "Turbot: turbot_saml_directory"
nav:
  title: turbot_saml_directory
---

# turbot\_saml\_directory

The `Turbot SAML Directory` resource adds support for creating SAML directories. It is used to create and delete directory settings.

## Example Usage

**Creating Your First SAML Directory**

```hcl
resource "turbot_saml_directory" "my_saml_directory" {
  parent                  = "tmod:@turbot/turbot#/"
  profile_id_template     = "myemail@myorg.com"
  description             = "This is our first SAML Directory"
  entry_point             = "https://example.com/myapp/sso/saml"
  certificate             = "sample-certificate"
}
```

## Argument Reference

The following arguments are supported:

- `parent` - (Required) The `id` or `aka` of the level at which the SAML directory will be created.
- `title` - (Required) Short descriptive name for the saml directory. This appears as the saml directory name in the Turbot Guardrails Console.
- `description` - (Optional) Brief description of the purpose and details of the directory.
- `entry_point` - (Required) Defines the identity provider single sign-on URL.
- `certificate` - (Required) The public key certificate ([base64-encoded](https://tools.ietf.org/html/rfc4648#section-4) ) which provides SAML entry point access
- `profile_id_template` - (Required) A template to generate profile id for users authenticated through this directory. For example, email id of the user.
- `issuer` - (Optional) a URL that uniquely identifies your SAML identity provider.
- `name_id_format` - (Optional) The name identifier format to request from the identity provider. Usually the Email Address is the accepted format.It accepts one of the two values - `UNSPECIFIED` and `EMAIL`. Defaults to `UNSPECIFIED`
- `sign_requests` - (Optional) Signing request for SAML authentication. It accepts one of the two values - `Enabled` and `Disabled`. If enabled, requests will be signed using the specified private key and signature algorithm.
- `signature_private_key` - (Optional) Private key used to sign authentication requests, in multiline PEM format starting with -----BEGIN PRIVATE KEY-----.
- `signature_algorithm` - (Optional) If a private key has been provided, it determines the signature algorithm for signing requests. If not specified defaults to *SHA-1*.
- `allow_group_syncing` -  (Optional) Boolean value to indicate whether groups will be synchronized for SAML users. Defaults to `false`.
- `allow_idp_initiated_sso` -  (Optional) Boolean value to indicate whether directory allows IDP-initiated SSO. Defaults to `false`.
- `profile_groups_attribute` - (Optional) Attribute returning list of groups that a SAML user is a part of.
- `group_filter` -  (Optional) Regular expression to filter out groups that are to be synced from SAML.
- `tags` - (Optional) User defined label for grouping resources.

## Attributes Reference

In addition to all the arguments above, the following attributes are PASSED :

- `id` - Unique identifier of the SAML directory.
- `parent_akas` - A list of all `akas` for the SAML directory's parent resource.
- `directory_type` - Type of the directory. For example, `saml`.
- `status` - Status of the SAML directory, which defaults to `Active`. Probable options are `Active`, `Inactive` and `New`.

## Import

SAML Directories can be imported using the `id`. For example,

```
terraform import turbot_saml_directory.my_saml_directory 123456789012
```
