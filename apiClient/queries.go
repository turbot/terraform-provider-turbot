package apiClient

import (
	"bytes"
	"fmt"
	"strings"
)

var turbotResourceMetadataFragment = `turbot {
  id
  parentId
  akas
  tags
}`

var turbotPolicyMetadataFragment = `turbot {
  id
  parentId
  akas
  tags
}`

func turbotGrantMetadataFragment(prefix string) string {
	return applyPrefix(prefix,
		`turbot {
  id
  profileId  
  resourceId 
}`)
}

func turbotActiveGrantMetadataFragment(prefix string) string {
	return applyPrefix(prefix, `turbot {
  id
  grantId
  resourceId
}`)
}

func applyPrefix(prefix, inputString string) string {
	return strings.Replace(inputString, "\n", "\n"+prefix, -1)
}

// validation
func validationQuery() (string, ValidationResponse) {
	query := `{
	schema: __schema {
		queryType {
			name
		}
	}
}`

	return query, ValidationResponse{}
}

// policySetting
func createPolicySettingMutation() string {
	return `mutation Create($command: PolicyCommandInput) {
	policySetting: policyCreate(command: $command ) {
		value: secretValue
		valueSource: secretValueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
		  id
		}
	}
}`

}

func readPolicySettingQuery(policySettingId string) string {
	return fmt.Sprintf(`{
	policySetting(id:"%s") {
		value: secretValue
		valueSource: secretValueSource
		template
        default
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
		  id
		}
	}
}`, policySettingId)
}

func updatePolicySettingMutation() string {
	return `mutation Update($command: PolicyCommandInput) {
	policySetting: policyUpdate(command: $command ) {
		value: secretValue
		valueSource: secretValueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
		  id
		}
	}
}`

}

func deletePolicySettingMutation() string {
	return `mutation Delete($command: PolicyCommandInput) {
	policySetting: policyDelete(command: $command ) {
		value: secretValue
		valueSource: secretValueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
		  id
		}
	}
}`
}

func findPolicySettingQuery(policyTypeUri, resourceAka string) string {
	return fmt.Sprintf(`{
  policySettings: policySettingList(filter: "policyType:%s resource:%s") {
    items {
      value
		valueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
		  id
		}
    }
  }
}
`, policyTypeUri, resourceAka)
}

// policy value
func readPolicyValueQuery(policyTypeUri string, resourceId string) string {
	return fmt.Sprintf(`{
	policyValue(uri:"%s", resourceId:"%s"){
		value: secretValue
		secretValue
		precedence
	    state
        reason
        details
        setting {
		  valueSource
          turbot {
            id
          }
        }
		turbot {
			id
		}
	}
}
`, policyTypeUri, resourceId)
}

// smart folder
func createSmartFolderMutation() string {
	return fmt.Sprintf(`mutation CreateSmartFolder($command: SmartFolderCommandInput) {
		smartFolder: smartFolderCreate(command: $command) {
			turbot {
				id
				parentId
				akas
			}
		}
	}`)
}

func readSmartFolderQuery(id string) string {
	return fmt.Sprintf(`{
	smartFolder: resource(id:"%s") {
		title: get(path:"title")
		description: get(path:"description")
		filter: get(path:"filters")
		parent:	get(path:"turbot.id")
		turbot: get(path:"turbot")
   		attachedResources{
			items{
				turbot{
					id
					akas
			}
		}
	}
 }
}`, id)
}

func updateSmartFolderMutation() string {
	return fmt.Sprintf(`mutation UpdateSmartFolder($command: SmartFolderCommandInput) {
		smartFolder: smartFolderUpdate(command: $command) {
			turbot {
				id
				parentId
				akas
			}
		}
	}`)
}

func createSmartFolderAttachmentMutation() string {
	return fmt.Sprintf(`mutation AttachSmartFolder($command: SmartFolderCommandInput) {
		smartFolderAttach(command: $command) {
			turbot {
				id
			}
		}
	}`)
}

func detachSmartFolderAttachment() string {
	return fmt.Sprintf(`mutation DetachSmartFolder($command: SmartFolderCommandInput) {
		detachSmartFolder: smartFolderDetach(command: $command) {
    		turbot {
				id
			}
  		}
	}`)
}

// mod
func installModMutation() string {
	return `mutation InstallMod($command: ModCommandInput) {
 	mod: modInstall(command: $command) {
		turbot {
		  id
          parentId
          akas
		}
		build
	}
}`
}

func readModQuery(modId string) string {
	return fmt.Sprintf(`{
  mod: resource(id:"%s") {
    uri: get(path: "turbot.akas.0")
    parent: get(path: "turbot.parentId")
    version: get(path: "version")
  }
}`, modId)
}

func uninstallModMutation() string {
	return `mutation UninstallMod($command: ModCommandInput) {
 	modUninstall(command: $command) {
		success
	}
}`
}

func modVersionsQuery(org, mod string) string {
	return fmt.Sprintf(`{
  versions: modVersionList(orgName: "%s", modName: "%s") {
    items {
      status
      version
    }
  }
}`, org, mod)
}

// resource
func createResourceMutation() string {
	return fmt.Sprintf(`mutation CreateResource($command: ResourceCommandInput) {
  resource: resourceCreate(command: $command) {
%s
  }
}`, turbotResourceMetadataFragment)
}

func readResourceQuery(aka string, properties map[string]string) string {
	var propertiesString bytes.Buffer
	if properties != nil {
		for alias, propertyPath := range properties {
			propertiesString.WriteString(fmt.Sprintf("    %s: get(path: \"%s\")\n", alias, propertyPath))
		}
	}
	return fmt.Sprintf(`{
  resource(id:"%s") {
%s
    turbot: get(path:"turbot")
  }
}`, aka, propertiesString.String())
}

func readGoogleDirectoryQuery(aka string) string {
	return fmt.Sprintf(`{
  directory: resource(id:"%s") {
	title:             get(path:"title")
	parent:            get(path:"turbot.parentId")
	description:       get(path:"description")
	status:            get(path:"status")
	directoryType:     get(path:"directoryType")
	profileIdTemplate: get(path:"profileIdTemplate")
	clientID:          get(path:"clientID")
	clientSecret:      getSecret(path:"clientSecret")
	poolId:            get(path:"poolId")
	groupIdTemplate:   get(path:"groupIdTemplate")
	loginNameTemplate: get(path:"loginNameTemplate")
	hostedName:        get(path:"hostedName")
    turbot: get(path:"turbot")
  }
}`, aka)
}

func updateResourceMutation() string {
	return `mutation UpsertResource($command: ResourceCommandInput) {
 	resource: resourceUpsert(command: $command) {
		turbot {
		  id
		  parentId
      akas
		}
	}
}`
}

func deleteResourceMutation() string {
	return `mutation DeleteResource($command: ResourceCommandInput) {
 	resource: resourceDelete(command: $command) {
		turbot {
		  id
		}
	}
}`
}

func readResourceListQuery(filter string, properties map[string]string) string {
	var propertiesString bytes.Buffer
	if properties != nil {
		for alias, propertyPath := range properties {
			propertiesString.WriteString(fmt.Sprintf("    %s: get(path: \"%s\")\n", alias, propertyPath))
		}
	}
	return fmt.Sprintf(`{
  resourceList(filter:"%s") {
	items{
		%s
    turbot: get(path:"turbot")
  	}
	}
}`, filter, propertiesString.String())
}

func readFullResourceQuery(aka string) string {
	return fmt.Sprintf(`{
  resource(id:"%s") {
    object
    turbot: get(path:"turbot")
  }
}`, aka)
}

// grant
func readGrantQuery(aka string) string {
	return fmt.Sprintf(`{
	grant: grant(id:"%s") {
		permissionTypeId
		permissionLevelId
		%s
	}
  }`, aka, turbotGrantMetadataFragment("\t\t"))
}

func createGrantMutation() string {
	return fmt.Sprintf(`mutation CreateGrant($command: GrantCommandInput) {
	grants: grantCreate(command: $command) {
		items{
			%s
		}
	}
}`, turbotGrantMetadataFragment("\t\t\t"))
}

func deleteGrantMutation() string {
	return fmt.Sprintf(`mutation DeleteGrant($command: GrantCommandInput) {
 	grant: grantDelete(command: $command) {
		items {
			%s
		}
	}
}`, turbotGrantMetadataFragment("\t\t\t"))
}

// active grant
func readActiveGrantQuery(aka string) string {
	return fmt.Sprintf(`{
	activeGrant: activeGrant(id:"%s"){
		%s
	}
}`, aka, turbotActiveGrantMetadataFragment("\t\t"))
}

func activateGrantMutation() string {
	return fmt.Sprintf(`mutation ActivateGrant($command: GrantCommandInput) {
	grantActivate: grantActivate(command: $command) {
		items {
			%s
		}
	}
}`, turbotActiveGrantMetadataFragment("\t\t\t"))
}

func deactivateGrantMutation() string {
	return fmt.Sprintf(`mutation DeactivateGrant($command: GrantCommandInput) {
	grantDeactivate(command: $command) {
		items {
			%s
		}
	}
}`, turbotActiveGrantMetadataFragment("\t\t\t"))
}
