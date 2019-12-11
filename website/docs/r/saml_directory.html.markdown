---
title: turbot_saml_directory
template: Documentation
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
  certificate             = "-----BEGIN CERTIFICATE-----MIICiTCCAfICCQD6m7oRw0uXOjANBgkqhkiG9w0BAQUFADCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wHhcNMTEwNDI1MjA0NTIxWhcNMTIwNDI0MjA0NTIxWjCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMaK0dn+a4GmWIWJ\n21uUSfwfEvySWtC2XADZ4nB+BLYgVIk60CpiwsZ3G93vUEIO3IyNoH/f0wYK8m9TrDHudUZg3qX4waLG5M43q7Wgc/MbQITxOUSQv7c7ugFFDzQGBzZswY6786m86gpE\nIbb3OhjZnzcvQAaRHhdlQWIMm2nrAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAtCu4nUhVVxYUntneD9+h8M\ng9q6q+auNKyExzyLwaxlAoo7TJHidbtS4J5iNmZgXL0Fkb\nFFBjvSfpJIlJ00zbhNYS5f6GuoEDmFJl0ZxBHjJnyp378OD8uTs7fLvjx79LjSTb\nNYiytVbZPQUQ5Yaxu2jXnimvw3rrszlaEXAMPLE=\n-----END CERIFICATE-----
}
```

## Argument Reference

The following arguments are supported:

- `parent` - (Required) The `id` or `aka` of the level at which the SAML directory will be created.
- `description` - (Required) Brief description of the purpose and details of the directory.
- `entry_point` - (Required) Defines the identity provider single sign-on URL.
- `certificate` - (Required) The public key certificate ( base64-encoded ) which provides SAML entry point access
- `profile_id_template` - (Required) A template to generate profile id for users authenticated through this directory. For example, email id of the user.
- `issuer` - (Optional) a URL that uniquely identifies your SAML identity provider.
- `group_id_template` - (Optional)  In case of a group profile, this template generates profile id for users authenticated through this directory. For example, email id of the group.
- `name_id_format` - (Optional) The name identifier format to request from the identity provider. Usually the Email Address is the accepted format.
- `sign_requests` - (Optional) Signing request for SAML authentication. It accepts one of the two values - `Enabled` and `Disabled`. If enabled, requests will be signed using the specified private key and signature algorithm.
- `signature_private_key` - (Optional) Private key used to sign authentication requests, in multiline PEM format starting with -----BEGIN PRIVATE KEY-----.
- `signature_algorithm` - (Optional) If a private key has been provided, it determines the signature algorithm for signing requests. If not specified defaults to *SHA-1*.
- `pool_id` - (Optional) Pool id associated with SAML directory.
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
