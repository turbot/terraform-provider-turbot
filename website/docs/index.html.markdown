---
title: Provider
template: Documentation
nav:
  order: 20
---

# Provider: turbot

The Turbot provider is used to interact with the resources supported by Turbot.
The provider needs to be configured with the proper credentials before it can
be used.

## Authentication

The Turbot provider offers a flexible means of providing credentials for authentication. The following methods are supported:

  - Credentials file
  - Static credentials
  - Environment variables

### Credentials file

The Turbot provider credentials can be authenticated using the Turbot credentials file. In this case you need to add the provider information in the Terraform configuration file. By default Turbot stores your `credentials.yml` file at a default location - `.config/turbot/`. To understand more about Turbot credentials, click [here](/docs/api/credentials).

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

  Static credentials can be provided by adding `access_key`, `secret_key` and `workspace` arguments in-line in the Turbot provider block. This information must be present in your configuration file.

**Example Usage**

  ```hcl
  # Configure the Turbot provider
  provider "turbot" {
    workspace           = "https://punisher-turbot.cloud.turbot-dev.com"
    access_key          = "b05*****-****-****-****-********580a"
    secret_key          = "d79*****-****-****-****-********b28"
  }

  # Create a new smart folder
  resource "turbot_smart_folder" "my_smart_folder" {
    # ...
  }

  # Set a policy
  resource "turbot_policy_setting" "s3_encryption_at_rest" {
    # ...
  }
  ```

### Environment Variables

You can provide your credentials via `TURBOT_ACCESS_KEY`, `TURBOT_SECRET_KEY` and `TURBOT_WORKSPACE` environment variables, representing your Turbot Access Key, Secret Key and workspace respectively.

**Example Usage**

```ruby
export TURBOT_SECRET_KEY=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export TURBOT_ACCESS_KEY=xxxxxxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export TURBOT_WORKSPACE=https://bananaman-turbot.putney.turbot.io
```

## Argument Reference

The following arguments are used:

* `workspace`  - Turbot workspace endpoint, e.g. `https://console-acme.cloud.turbot.com/api/latest/graphql`. May also be set via the `TURBOT_WORKSPACE` environment variable.
* `access_key` - Turbot access key, e.g. `c32ee14d-615b-4efb-95c3-0cf3f680d2fc`. May also be set via the `TURBOT_ACCESS_KEY` environment variable.
* `secret_key` - Turbot secret key, e.g. `a2d6660d-0feb-42c7-9718-274cb5a82ed7`. May also be set via the `TURBOT_SECRET_KEY` environment variable.
* `profile`    - Turbot workspace profile, e.g. `testProfile`. May also be set via the `TURBOT_PROFILE` environment variable.
* `credentials_file`    - Turbot shared credentials path, e.g. `user/testUser/{{credential_file_path}}`. May also be set via the `TURBOT_SHARED_CREDENTIALS_PATH` environment variable.