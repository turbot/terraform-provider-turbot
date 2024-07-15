---
layout: "turbot"
title: Provider
template: Documentation
nav:
  order: 20
---

# Turbot Guardrails Provider

The Turbot Guardrails provider is used to interact with the resources supported by Turbot Guardrails.
The provider needs to be configured with the proper credentials before it can
be used.

## Authentication

The Turbot Guardrails provider offers a flexible means of providing credentials for authentication. The following methods are supported:

  - Credentials file
  - Static credentials
  - Environment variables

### Credentials file

The Turbot Guardrails provider credentials can be authenticated using the Guardrails credentials file. In this case you need to add the provider information in the Terraform configuration file. By default Guardrails stores your `credentials.yml` file at a default location - `.config/turbot/`.

**Example (Using your default profile)**

   ```hcl
     provider "turbot" {}
   ```

**Example (Using a named profile)**

   ```hcl
    provider "turbot" {
        profile  = "MyProfile"
    }
   ```

Alternatively you can also store your credentials in your desired path. This can be accessed by using the `credentials_file` argument along with `profile` argument.

**Example Usage**

  ```hcl
   provider "turbot" {
     profile                  = MyProfile
     credentials_file         = "/Users/test_user_name/{{credential_file_path}}"
   }
  ```

### Static Credentials

  Static credentials can be provided by adding `access_key`, `secret_key` and `workspace` arguments in-line in the Turbot Guardrails provider block. This information must be present in your configuration file.

**Example Usage**

  ```hcl
  # Configure the Turbot Guardrails provider
  provider "turbot" {
    workspace           = "https://example.com"
    access_key          = "b05*****-****-****-****-********580a"
    secret_key          = "d79*****-****-****-****-********b28"
  }

  # Create a new policy pack
  resource "turbot_policy_pack" "my_policy_pack" {
    # ...
  }

  # Set a policy
  resource "turbot_policy_setting" "s3_encryption_at_rest" {
    # ...
  }
  ```

### Environment Variables

You can provide your credentials via `TURBOT_ACCESS_KEY`, `TURBOT_SECRET_KEY` and `TURBOT_WORKSPACE` environment variables, representing your Turbot Guardrails Access Key, Secret Key and workspace respectively.

**Example Usage**

   ```ruby
    export TURBOT_SECRET_KEY=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    export TURBOT_ACCESS_KEY=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    export TURBOT_WORKSPACE=https://example.com
   ```

## Argument Reference

The following arguments are used:

* `workspace`  - Turbot Guardrails workspace endpoint, e.g. `https://example.com/api/latest/graphql`. May also be set via the `TURBOT_WORKSPACE` environment variable.
* `access_key` - Turbot Guardrails access key, e.g. `1wxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxe6`. May also be set via the `TURBOT_ACCESS_KEY` environment variable.
* `secret_key` - Turbot Guardrails secret key, e.g. `b90xxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxnp`. May also be set via the `TURBOT_SECRET_KEY` environment variable.
* `profile`    - Turbot Guardrails workspace profile, e.g. `testProfile`. May also be set via the `TURBOT_PROFILE` environment variable.
* `credentials_file`    - Turbot Guardrails shared credentials path, e.g. `user/testUser/{{credential_file_path}}`. May also be set via the `TURBOT_SHARED_CREDENTIALS_PATH` environment variable.
