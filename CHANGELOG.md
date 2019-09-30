

## 1.0.0-beta.4 (September 30, 2019)

FEATURES:
* **New Resource:** `turbot_grant` ([#31](https://github.com/turbotio/terraform-provider-turbot/issues/31))
* **New Resource:** `turbot_smart_folder` ([#39](https://github.com/turbotio/terraform-provider-turbot/issues/39))
* **New Resource:** `turbot_smart_folder_attachment` ([#39](https://github.com/turbotio/terraform-provider-turbot/issues/39))
* **New Resource:** `turbot_shadow_resource` ([#62](https://github.com/turbotio/terraform-provider-turbot/issues/62))

ENHANCEMENTS:
* Add tags support to various resources  ([#55](https://github.com/turbotio/terraform-provider-turbot/issues/55)): 
  * turbot_folder
  * turbot_resource 
  * turbot_local_directory
  * turbot_saml_directory
  * turbot_google_directory
  

## 1.0.0-beta.3 (September 20, 2019)

FEATURES:
* **New Resource:** `turbot_saml_directory` ([#34](https://github.com/turbotio/terraform-provider-turbot/issues/34))
* **New Resource:** `turbot_google_directory` ([#39](https://github.com/turbotio/terraform-provider-turbot/issues/39))

## 1.0.0-beta.2 (September 18, 2019)

FEATURES:
* **New Data Source:** `turbot_resource` ([#20](https://github.com/turbotio/terraform-provider-turbot/issues/20))
* **New Resource:** `turbot_local_directory` ([#14](https://github.com/turbotio/terraform-provider-turbot/issues/14))
* **New Resource:** `turbot_local_directory_user` ([#26](hhttps://github.com/turbotio/terraform-provider-turbot/issues/26))
* **New Resource:** `turbot_profile` ([#16](https://github.com/turbotio/terraform-provider-turbot/issues/16))

BUGFIXES:

* When searching for existing policy setting before creation, ignore default setting. ([#9](https://github.com/turbotio/terraform-provider-turbot/issues/9))

ENHANCEMENTS:

* resource/turbot_folder resource: remove error when creating folder with existing name and parent. This is permitted. ([#12](https://github.com/turbotio/terraform-provider-turbot/issues/12))
* resource/turbot_policy_setting: add default value of "required" for precedence property.  ([#13](https://github.com/turbotio/terraform-provider-turbot/issues/13))

## 1.0.0-beta.1 (September 13, 2019)

ENHANCEMENTS:

* resource/turbot_mod: support version ranges, plan shows changes if a new version has been published ([#6](https://github.com/turbotio/terraform-provider-turbot/issues/6))

FEATURES:

* **New Resource:** `turbot_resource` ([#7](https://github.com/turbotio/terraform-provider-turbot/issues/7))


## 1.0.0-beta.0 (September 12, 2019)

FEATURES:

* **New Resource:** `turbot_mod` ([#1](https://github.com/turbotio/terraform-provider-turbot/issues/1))
* **New Resource:** `turbot_folder` ([#2](https://github.com/turbotio/terraform-provider-turbot/issues/2))
* **New Resource:** `turbot_policy_setting` ([#3](https://github.com/turbotio/terraform-provider-turbot/issues/3))
* **New Data Source:** `turbot_policy_value` ([#4](https://github.com/turbotio/terraform-provider-turbot/issues/4))
