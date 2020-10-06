## 1.7.0 (October 5, 2020)
FEATURES:
* **New Resource:** `turbot_ldap_directory` ([#56](https://github.com/turbot/terraform-provider-turbot/issues/56))

ENHANCEMENTS:
* `resource/resource_turbot_resource`: Introduced arguments `full_data` and `full_metadata` to provide alternatives to `data` and `metadata`. These allow you to fully manage the resource data/metadata(any additional properties not in the config are deleted). Refer to resource documentation for full details.

BUG FIXES:
* `resource/resource_turbot_file`: When creating a resource with no `content` attribute set, `apply` operation fails. ([#117](https://github.com/turbot/terraform-provider-turbot/issues/117))
* `resource/resource_turbot_folder`: After deletion of `description` attribute, `plan` operation shows a diff. ([#127](https://github.com/turbot/terraform-provider-turbot/issues/127))
* The profile argument in provider config had lower precedence than TURBOT_PROFILE set in environment variables. This has been corrected to give the provider config higher precedence.

## 1.6.3 (August 24, 2020)
BUG FIXES:
* `resource/resource_turbot_mod`: After successful import, an additional mod install operation is triggered when running terraform apply. ([#95](https://github.com/turbot/terraform-provider-turbot/issues/75))
* Certain errors caused the error handling code to crash. ([#108](https://github.com/turbot/terraform-provider-turbot/issues/108))

## 1.6.2 (August 17, 2020)
BUG FIXES:
* `resource/resource_turbot_google_directory`: After creation of resource, `terraform plan` shows a change required for the attribute `hosted_name`. ([#75](https://github.com/turbot/terraform-provider-turbot/issues/75))
* `resource/resource_turbot_profile`: After creation of resource, `terraform plan` shows a change required for the attribute `picture`. ([#78](https://github.com/turbot/terraform-provider-turbot/issues/78))

## 1.6.1 (July 28, 2020)
Transfer provider to Hashicorp registry. No user-facing changes.

## 1.6.0 (July 20, 2020)
FEATURES:
* **New Data Source:** `turbot_control` ([#32](https://github.com/terraform-providers/terraform-provider-turbot/issues/32))

ENHANCEMENTS:
* `resource/resource_turbot_saml_directory`:  Add argument `allow_idp_initiated_sso`to provide support for Allow IdP-initiated SSO. ([#39](https://github.com/terraform-providers/terraform-provider-turbot/issues/39))
* `resource/resource_turbot_mod`: Add `timeout` block to support `create` timeout. This is the length of time  provider waits for mod installation to complete. Default value is 15min ([#53](https://github.com/terraform-providers/terraform-provider-turbot/issues/53))
## 1.5.0 (July 15, 2020)
ENHANCEMENTS:
* Internal errors (5xx) received from the Turbot server will now show a friendly error message. We hope you never see them, but if you do, they should be easier to understand. ([#54](https://github.com/terraform-providers/terraform-provider-turbot/issues/54))

## 1.4.0 (July 09, 2020)
FEATURES:
* **New Resource:** `turbot_file` ([#34](https://github.com/terraform-providers/terraform-provider-turbot/issues/34))

BUG FIXES:
* `resource/resource_turbot_saml_directory`: Deprecate `group_id_template` and `pool_id` arguments, as these are no longer required or supported by the Turbot API. If passed, these parameters will be ignored. ([#37](https://github.com/terraform-providers/terraform-provider-turbot/issues/37))
* `resource/resource_turbot_local_directory`: Deprecate `group_id_template`, `login_name_template` and `pool_id` arguments, as these are no longer required or supported by the Turbot API. If passed, these parameters will be ignored. ([#37](https://github.com/terraform-providers/terraform-provider-turbot/issues/37))
## 1.3.1 (July 02, 2020)
BUG FIXES 
* `resource/resource_turbot_mod`: After creation of resource, when running `terraform plan`, `version_current` attribute gives a `diff` .([#38](https://github.com/terraform-providers/terraform-provider-turbot/issues/38))

## 1.3.0 (May 29, 2020)
ENHANCEMENTS:
* `resource/resource_turbot_policy_setting`: `template_input` attribute now supports either a string or an array of strings in `YAML` format. ([#11](https://github.com/terraform-providers/terraform-provider-turbot/issues/11))
* `resource/resource_turbot_shadow_resource`: Add `timeout` block to support `create` timeout. Default value of `create` timeout is 5min ([#19](https://github.com/terraform-providers/terraform-provider-turbot/issues/19))

BUG FIXES 
* `resource/resource_turbot_smart_folder_attachment`: After importing a resource, when running `terraform plan`, a resource replacement would be required.([#8](https://github.com/terraform-providers/terraform-provider-turbot/issues/8))

## 1.2.0 (May 20, 2020)
FEATURES:
* **New Resource:** `turbot_turbot_directory` ([#1](https://github.com/terraform-providers/terraform-provider-turbot/issues/1))

ENHANCEMENTS:
* `resource/resource_turbot_saml_directory`:  Add arguments `allow_group_syncing`,`profile_groups_attribute`, `group_filter` and default value of `name_id_format` to `UNSPECIFIED`. ([#3](https://github.com/terraform-providers/terraform-provider-turbot/issues/3))

BUG FIXES
* `resource/resource_turbot_resource`: Fix `Data validation error` while doing a update on resource.([#13](https://github.com/terraform-providers/terraform-provider-turbot/issues/13))

## 1.1.0 (April 21, 2020)
BUG FIXES
* `resource/resource_turbot_google_directory`: For security, the client secret is not read from the resource. This causes diffs to be identified after importing the resource, requiring replacement of the resource. Suppress diffs caused by client_secret.([#66](https://github.com/turbot/terraform-provider-turbot/issues/66))

## 1.0.4 (February 26, 2020)
BUG FIXES
* `resource/resource_turbot_resource`: Import was not setting the type property, causing terraform plan to require replacement of the resource.([#62](https://github.com/turbot/terraform-provider-turbot/issues/62))

## 1.0.3 (February 25, 2020)
BUG FIXES
* `resource/resource_turbot_policy_setting`: Import was not setting the type property, causing terraform plan to require replacement of the resource.([#57](https://github.com/turbot/terraform-provider-turbot/issues/57))

TECHNICAL:
* `resource/resource_turbot_google_directory`:  Invalid properties were being passed to the turbot API.([#56](https://github.com/turbot/terraform-provider-turbot/issues/56))

## 1.0.2 (February 11, 2020)
BUG FIXES
* `resource/resource_turbot_resource`: Import was failing with error: `Unexpected JSON`.([#50](https://github.com/turbot/terraform-provider-turbot/issues/50))

## 1.0.1 (February 5, 2020)
ENHANCEMENTS:
* `resource/resource_turbot_shadow_resource`: Improved error handling with retries to improve resilience.([#42](https://github.com/turbot/terraform-provider-turbot/issues/42))

## 1.0.0 (December 18, 2019)
GENERAL
* Add MPL license

## 1.0.0-beta.13 (December 12, 2019)
BUG FIXES
* `resource/resource_turbot_smart_folder`: Smart folder creation was failing due to invalid graphql mutation.([#183](https://github.com/turbotio/terraform-provider-turbot/issues/183))

TECHNICAL:
* Update acceptance test function names to have underscores.([#186](https://github.com/turbotio/terraform-provider-turbot/issues/186))

ENHANCEMENTS:
* Improved error logging. ([#182](https://github.com/turbotio/terraform-provider-turbot/issues/182))

## 1.0.0-beta.12 (December 12, 2019)
BUG FIXES
* `resource/turbot_shadow_resource`: When an aka is provided in the `resource` argument, multiple results are sometimes being returned. Use a `resource` query instead of a `resourceList` query with the filter `resource:<aka>`.([#176](https://github.com/turbotio/terraform-provider-turbot/issues/176))

ENHANCEMENTS:
* Add documentation for resources and data sources ([#178](https://github.com/turbotio/terraform-provider-turbot/issues/178))
* Add vendor folder to git ([#181](https://github.com/turbotio/terraform-provider-turbot/issues/181))

## 1.0.0-beta.11 (December 6, 2019)
BUG FIXES
* `resource/turbot_shadow_resource`: `read error`. When doing the read of tracked resource, the shadow resource read query returns `child elements` as well. Fix this return tracked resource.([#162](https://github.com/turbotio/terraform-provider-turbot/issues/162))  

TECHNICAL:
* Fix `buildProperties` function to support `implicit & explicit` mapping ([#164](https://github.com/turbotio/terraform-provider-turbot/issues/164))
* `resource/turbot_folder` : update `create/update` mutations to accomodate `read` operation attribute changes([#154](https://github.com/turbotio/terraform-provider-turbot/issues/154))
* `resource/turbot_local_directory` : update `create/update` mutations to accomodate `read` operation attribute changes([#156](https://github.com/turbotio/terraform-provider-turbot/issues/156))
* `resource/turbot_local_directory_user` : update `create/update` mutations to accomodate `read` operation attribute changes([#160](https://github.com/turbotio/terraform-provider-turbot/issues/160))
* `resource/turbot_saml_directory` : update `create/update` mutations to accomodate `read` operation attribute changes([#166](https://github.com/turbotio/terraform-provider-turbot/issues/166))
* `resource/turbot_mod` : update `create/update` mutations to accomodate `read` operation attribute changes([#170](https://github.com/turbotio/terraform-provider-turbot/issues/170))
* `resource/turbot_grant_activation` : update `create/update` mutations to accomodate `read` operation attribute changes([#168](https://github.com/turbotio/terraform-provider-turbot/issues/168))
* `resource/turbot_policy_setting` : update `create/update` mutations to accomodate `read` operation attribute changes([#168](https://github.com/turbotio/terraform-provider-turbot/issues/168))
* `resource/turbot_google_directory` : update `create/update` mutations to accomodate `read` operation attribute changes([#168](https://github.com/turbotio/terraform-provider-turbot/issues/168))
* `resource/turbot_smart_folder` : update `create/update` mutations to accomodate `read` operation attribute changes([#168](https://github.com/turbotio/terraform-provider-turbot/issues/168))
* `resource/turbot_profile` : update `create/update` mutations to accomodate `read` operation attribute changes([#158](https://github.com/turbotio/terraform-provider-turbot/issues/154))

## 1.0.0-beta.10 (November 20, 2019)

BUG FIXES
* Update Error handling for `not found` and `data validation` errors to correctly match error strings([#143](https://github.com/turbotio/terraform-provider-turbot/issues/143))  

## 1.0.0-beta.9 (November 11, 2019)

BUG FIXES
* `resource/turbot_shadow_resource` : `deletion error`. When the tracked resource is deleted, the shadow resource would display an error. Fix this to delete the shadow resource when the tracked resource is deleted.([#127](https://github.com/turbotio/terraform-provider-turbot/issues/127)) 

ENHANCEMENTS:
* `resource/turbot_profile` : make `directory_pool_id` optional
* `resource/turbot_folder` : make `description` optional
## 1.0.0-beta.8 (November 4, 2019)
NOTES: 
This requires core version >=5.0.0-beta.96  

ENHANCEMENTS:
* `resource/turbot_policy_setting` : change default precedence to `REQUIRED`. ([#125](https://github.com/turbotio/terraform-provider-turbot/issues/125))
* `resource/turbot_policy_setting` : add `resource_akas` property and diff suppression. ([#116](https://github.com/turbotio/terraform-provider-turbot/issues/116))

BUG FIXES

* `data-source/turbot_policy_value` : handle resource not found error. ([#124](https://github.com/turbotio/terraform-provider-turbot/issues/124))

## 1.0.0-beta.7 (October 22, 2019)

BREAKING CHANGES
* `resource/turbot_grant` : change `profile` property to `identity` ([#106](https://github.com/turbotio/terraform-provider-turbot/issues/106)) 

TECHNICAL:
* `resource/turbot_grant` : to use new mutations ([#109](https://github.com/turbotio/terraform-provider-turbot/issues/109))
* `resource/turbot_grant_activation` : to use new mutations ([#110](https://github.com/turbotio/terraform-provider-turbot/issues/110))
* `resource/turbot_mod` : to use new mutations ([#111](https://github.com/turbotio/terraform-provider-turbot/issues/111))
* `resource/turbot_policy_setting` : to use new mutations ([#112](https://github.com/turbotio/terraform-provider-turbot/issues/112))
* `resource/turbot_smart_folder` : to use new mutations ([#104](https://github.com/turbotio/terraform-provider-turbot/issues/104))
* `resource/turbot_smart_folder_attachment` : to use new mutations ([#104](https://github.com/turbotio/terraform-provider-turbot/issues/104))

## 1.0.0-beta.6 (October 14, 2019)

ENHANCEMENTS:
* Add support for credentials profiles. ([#57](https://github.com/turbotio/terraform-provider-turbot/issues/57))
* `resource/turbot_grant` : change arguments `permission_type` and `permission_level` to `type` and `level`. ([#92](https://github.com/turbotio/terraform-provider-turbot/issues/92))
* `resource/turbot_policy_setting` : rename `policy_type` to `type`. ([#87](https://github.com/turbotio/terraform-provider-turbot/issues/87)) 
* `resource/turbot_policy_value` : rename `policy_type` to `type`. ([#86](https://github.com/turbotio/terraform-provider-turbot/issues/86)) 
* `resource/turbot_resource` : data provider schema. ([#42](https://github.com/turbotio/terraform-provider-turbot/issues/42))  
* `resource/turbot_mod` : default `parent` property to the turbot resource. ([#93](https://github.com/turbotio/terraform-provider-turbot/issues/93)) 
* Update client to support renamed credentials environment variables `TURBOT_ACCESS_KEY` and `TURBOT_SECRET_KEY` ([#90](https://github.com/turbotio/terraform-provider-turbot/issues/90))

TECHNICAL
* Move MapFromResourceData and StoreAkas to helpers. ([#89](https://github.com/turbotio/terraform-provider-turbot/issues/89)) 
* Update all resources to support new resource mutation schema.  ([#94](https://github.com/turbotio/terraform-provider-turbot/issues/94))
* Create helpers package. ([#89](https://github.com/turbotio/terraform-provider-turbot/issues/89))
  

BUG FIXES
* Update client.BuildApiUrl to require both workspace and installation domain to be provided (e.g. `bananaman-turbot.putney.turbot.io`, rather than just `bananaman-turbot.putney`) ([#98](https://github.com/turbotio/terraform-provider-turbot/issues/98))

## 1.0.0-beta.5 (October 03, 2019)

BUG FIXES: 
* `resource/turbot_mod` : failing to install mods. Add error handling to the code to check for existing mod. ([#82](https://github.com/turbotio/terraform-provider-turbot/issues/82))

FEATURES:
* **New Resource:** `turbot_grant_activation` ([#79](https://github.com/turbotio/terraform-provider-turbot/issues/79))

ENHANCEMENTS:
* Add support for terraform 0.12. ([#75](https://github.com/turbotio/terraform-provider-turbot/issues/75))
* Update all directory resource schemas to make status and directory_type computed.([#76](https://github.com/turbotio/terraform-provider-turbot/issues/76)) 
* `resource/turbot_policy_setting` : support encryption of value and value_source in state file. ([#77](https://github.com/turbotio/terraform-provider-turbot/issues/77))
* `resource/turbot_google_directory` : support encryption of client_secret in state file. ([#47](https://github.com/turbotio/terraform-provider-turbot/issues/47))

## 1.0.0-beta.4 (September 30, 2019)

FEATURES:
* **New Resource:** `turbot_grant` ([#31](https://github.com/turbotio/terraform-provider-turbot/issues/31))
* **New Resource:** `turbot_smart_folder` ([#39](https://github.com/turbotio/terraform-provider-turbot/issues/39))
* **New Resource:** `turbot_smart_folder_attachment` ([#39](https://github.com/turbotio/terraform-provider-turbot/issues/39))
* **New Resource:** `turbot_shadow_resource` ([#62](https://github.com/turbotio/terraform-provider-turbot/issues/62))

ENHANCEMENTS:
* Add tags support to various resources  ([#55](https://github.com/turbotio/terraform-provider-turbot/issues/55)): 
  * `turbot_folder`
  * `turbot_resource `
  * `turbot_local_directory`
  * `turbot_saml_directory`
  * `turbot_google_directory`
  

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

BUG FIXES:

* resource/turbot_policy_setting: when searching for existing policy setting before creation, ignore default setting. ([#9](https://github.com/turbotio/terraform-provider-turbot/issues/9))

ENHANCEMENTS:

* `resource/turbot_folder` : remove error when creating folder with existing name and parent. This is permitted. ([#12](https://github.com/turbotio/terraform-provider-turbot/issues/12))
* `resource/turbot_policy_setting` : add default value of "required" for precedence property.  ([#13](https://github.com/turbotio/terraform-provider-turbot/issues/13))

## 1.0.0-beta.1 (September 13, 2019)

ENHANCEMENTS:

* `resource/turbot_mod` : support version ranges, plan shows changes if a new version has been published ([#6](https://github.com/turbotio/terraform-provider-turbot/issues/6))

FEATURES:

* **New Resource:** `turbot_resource` ([#7](https://github.com/turbotio/terraform-provider-turbot/issues/7))


## 1.0.0-beta.0 (September 12, 2019)

FEATURES:

* **New Resource:** `turbot_mod` ([#1](https://github.com/turbotio/terraform-provider-turbot/issues/1))
* **New Resource:** `turbot_folder` ([#2](https://github.com/turbotio/terraform-provider-turbot/issues/2))
* **New Resource:** `turbot_policy_setting` ([#3](https://github.com/turbotio/terraform-provider-turbot/issues/3))
* **New Data Source:** `turbot_policy_value` ([#4](https://github.com/turbotio/terraform-provider-turbot/issues/4))
