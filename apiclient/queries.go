package apiclient

import (
	"bytes"
	"fmt"
)

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
		}
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
	return `mutation CreateResource($command: ResourceCommandInput) {
 	resource: resourceCreate(command: $command) {
		turbot {
		  id
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
