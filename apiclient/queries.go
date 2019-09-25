package apiclient

import (
	"bytes"
	"fmt"
)

var turbotMetadataQuery = `turbot {
  id
  parentId
  akas
  tags
}`

// return query and matching response object to receive query result

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

// create policySetting
func createPolicySettingMutation() string {
	return `mutation Create($command: PolicyCommandInput) {
	policySetting: policyCreate(command: $command ) {
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
}`

}

// read policySetting
func readPolicySettingQuery(policySettingId string) string {
	return fmt.Sprintf(`{
	policySetting(id:"%s") {
		value
		valueSource
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

// update policySetting
func updatePolicySettingMutation() string {
	return `mutation Update($command: PolicyCommandInput) {
	policySetting: policyUpdate(command: $command ) {
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
}`

}

// delete policySetting
func deletePolicySettingMutation() string {
	return `mutation Delete($command: PolicyCommandInput) {
	policySetting: policyDelete(command: $command ) {
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
}`
}

// create smart folder mutation
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

//update smart folder mutation

func updateSmartFolderMutation() string {
	return fmt.Sprintf(`mutation UpdateSmartFolder($command: SmartFolderCommandInput) {
		smartFolder:	smartFolderUpdate(command: $command) {
			turbot {
				id
				parentId
				akas
			}
		}
	}`)
}

// attach folderResource Mutation

func createSmartFolderAttachmentMutation() string {
	return fmt.Sprintf(`mutation AttachSmartFolder($command: SmartFolderCommandInput) {
		smartFolder :smartFolderAttach(command: $command) {
			turbot {
				id
			}
		}
	}`)
}

// detach FolderResource mutation

func detachSmartFolderAttachment() string {
	return fmt.Sprintf(`mutation DetachSmartFolder($command: SmartFolderCommandInput) {
		detachSmartFolder: smartFolderDetach(command: $command) {
    		turbot {
				id
			}
  		}
	}`)
}

// read Smart Folder attached to multiple Resources

func readAttachedResourcesOnSmartfolder(id string) string {
	return fmt.Sprintf(`
	resource(id:"%s") {
   	attachedResources{
			items{
				
			}
		}
 }
`, id)
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

// read policyValue
func readPolicyValueQuery(policyTypeUri string, resourceId string) string {
	return fmt.Sprintf(`{
	policyValue(uri:"%s", resourceId:"%s"){
		value
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

// install mod
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

// uninstall mod
func uninstallModMutation() string {
	return `mutation UninstallMod($command: ModCommandInput) {
 	modUninstall(command: $command) {
		success
	}
}`
}

// get mod versions
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

// read mod
func readModQuery(modId string) string {
	return fmt.Sprintf(`{
  mod: resource(id:"%s") {
    uri: get(path: "turbot.akas.0")
    parent: get(path: "turbot.parentId")
    version: get(path: "version")
  }
}`, modId)
}

// create resource
func createResourceMutation() string {
	return fmt.Sprintf(`mutation CreateResource($command: ResourceCommandInput) {
  resource: resourceCreate(command: $command) {
%s
  }
}`, turbotMetadataQuery)
}

// create grant
func createGrantMutation() string {
	return `mutation CreateGrant($command: GrantCommandInput) {
	grants: grantCreate(command: $command) {
		items{
			turbot{
				id
				resourceId
			}
		}
	}
}`
}

// update resource
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

// delete resource
func deleteResourceMutation() string {
	return `mutation DeleteResource($command: ResourceCommandInput) {
 	resource: resourceDelete(command: $command) {
		turbot {
		  id
		}
	}
}`
}

// read resource
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

func readGrantQuery(aka string) string {
	return fmt.Sprintf(`{
		grant: grant(id:"%s") {
			permissionTypeId
			permissionLevelId
				turbot{
                    id
					resourceId,
					profileId
				}
		}
	  }`, aka)
}

// delete resource
func deleteGrantMutation() string {
	return `mutation DeleteGrant($command: GrantCommandInput) {
 	grant: grantDelete(command: $command) {
		items {
			turbot {
		  		id
			}
		}
	}
}`
}
