

## 1.0.0-beta.6 (tbd)

ENHANCEMENTS:
Update turbot_resource data provider schema. ([#42](https://github.com/turbotio/terraform-provider-turbot/issues/42))  
Create helpers package. ([#89](https://github.com/turbotio/terraform-provider-turbot/issues/89))
Add support for credentials profiles. ([#57](https://github.com/turbotio/terraform-provider-turbot/issues/57))
Update turbot_policy_setting - rename policy_type to type. ([#87](https://github.com/turbotio/terraform-provider-turbot/issues/87)) 
Update turbot_policy_value - rename policy_type to type. ([#86](https://github.com/turbotio/terraform-provider-turbot/issues/86)) 

## 1.0.0-beta.5 (October 03, 2019)

BUGFIXES: 
* turbot_mod failing to install mods. Add error handling to the code to check for existing mod. ([#82](https://github.com/turbotio/terraform-provider-turbot/issues/82))

FEATURES:
* **New Resource:** `turbot_grant_activation` ([#79](https://github.com/turbotio/terraform-provider-turbot/issues/79))

ENHANCEMENTS:
* Add support for terraform 0.12. ([#75](https://github.com/turbotio/terraform-provider-turbot/issues/75))
* Update all directory resource schemas to make status and directory_type computed.([#76](https://github.com/turbotio/terraform-provider-turbot/issues/76)) 
* Update turbot_policy_setting - support encryption of value and value_source in state file. ([#77](https://github.com/turbotio/terraform-provider-turbot/issues/77))
* Update turbot_google_directory - support encryption of client_secret in state file. ([#47](https://github.com/turbotio/terraform-provider-turbot/issues/47))

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
