package apiClient

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/blang/semver"
)

// NOTE: we do not use a fragment for resource metadata as we just request the full turbot property
// using turbot: get(path:"turbot")
// This is because we saw errors returning null for the turbot property for a non existent resource
// TODO fix this to use a fragment

func turbotPolicyMetadataFragment(prefix string) string {
	return applyPrefix(prefix, `turbot {
	id
	parentId
	akas
	tags
}`)
}

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

// add the given prefix to each line of the multi-line inputString
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
	return `mutation CreatePolicySetting($input: CreatePolicySettingInput!) {
	policySetting: createPolicySetting(input: $input ) {
		type {
			uri
		}
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

// The id is passed as a GraphQL variable, never interpolated into the document. Interpolating a
// caller-supplied identifier into a string literal lets a value containing a double quote escape
// the literal and append attacker-chosen GraphQL, executed with the provider's credentials (see
// TestReadBuildersDoNotInterpolateIdentifier). policySetting.id is a nullable ID; a non-null
// variable is assignable to it.
func readPolicySettingQuery() string {
	return `query ReadPolicySetting($id: ID!) {
policySetting(id: $id) {
	type {
		uri
	}
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
		resourceId
	}
}
}`
}

func updatePolicySettingMutation() string {
	return `mutation UpdatePolicySetting($input: UpdatePolicySettingInput!) {
	policySetting: updatePolicySetting(input: $input ) {
		type {
			uri
		}
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
	return `mutation DeletePolicySetting($input: DeletePolicySettingInput!) {
	policySetting: deletePolicySetting(input: $input ) {
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

func findPolicyTypeQuery(policyTypeUri string) string {
	return fmt.Sprintf(`{
  policyTypes: policyTypes(filter: "policyTypeId:%s level:self") {
    items {
		modUri
		turbot {
			id
		}
    }
  }
}
`, policyTypeUri)
}

// watch
func createWatchMutation() string {
	return fmt.Sprintf(`mutation CreateWatch($input: CreateWatchInput!) {
		watch: createWatch(input: $input) {
			description
			filters
			handler
			turbot {
				id
				resourceId
				favoriteId
			}
		}
	}`)
}

func readWatchQuery() string {
	return `query ReadWatch($id: ID!) {
		watch(id: $id) {
			description
			filters
			handler
			turbot {
				id
				resourceId
				favoriteId
			}
		}
	}`
}

func updateWatchMutation() string {
	return fmt.Sprintf(`mutation UpdateWatch($input: UpdateWatchInput!) {
		updateWatch(input: $input) {
			description
			filters
			handler
			turbot {
				id
				resourceId
				favoriteId
			}
		}
	}`)
}

func deleteWatchMutation() string {
	return `mutation DeleteMyWatch($id: ID!) {
		deleteWatch(input: {id: $id}) {
			handler
			filters
		}
	}
	`
}

// smart folder
// filter and description are removed for a workaround, will be removed after a Core change.
func createSmartFolderMutation() string {
	return `mutation CreateSmartFolder($input: CreateSmartFolderInput!) {
		smartFolder: createSmartFolder(input: $input) {
			turbot {
				id
				parentId
				akas
				title
			}
		}
	}`
}

func readSmartFolderQuery() string {
	return `query ReadSmartFolder($id: ID!) {
	smartFolder: resource(id: $id) {
		title: get(path:"turbot.title")
		description: get(path:"description")
		filters: get(path:"filters")
		parent:	get(path:"turbot.parentId")
		turbot: get(path:"turbot")
	}
}`
}

func updateSmartFolderMutation() string {
	return `mutation UpdateSmartFolder($input: UpdateSmartFolderInput!) {
		smartFolder: updateSmartFolder(input: $input) {
			turbot {
				id
				parentId
				akas
			}
		}
	}`
}

func deleteSmartFolderMutation() string {
	return `mutation DeleteSmartFolder($id: ID!) {
		smartFolder: deleteSmartFolder(input: {id: $id}) {
			turbot {
				id
			}
		}
	}`
}

func createSmartFolderAttachmentMutation() string {
	return `mutation AttachSmartFolder($input: AttachSmartFolderInput!) {
		attachSmartFolders(input: $input) {
			turbot {
				id
			}
		}
	}`
}

func detachSmartFolderAttachment() string {
	return `mutation DetachSmartFolder($input: DetachSmartFolderInput!) {
		detachSmartFolder: detachSmartFolders(input: $input) {
    		turbot {
				id
			}
  		}
	}`
}

// mod
func installModMutation() string {
	return `mutation InstallMod($input: InstallModInput!) {
	mod: installMod(input: $input) {
		turbot {
			id
			parentId
			akas
		}
		build
	}
}`
}

func readModQuery() string {
	return `query ReadMod($id: ID!) {
	mod: resource(id: $id) {
		uri: get(path: "turbot.akas.0")
		parent: get(path: "turbot.parentId")
		version: get(path: "version")
	}
}`
}

func uninstallModMutation() string {
	return `mutation UninstallMod($input: UninstallModInput!) {
	uninstallMod(input: $input) {
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
func createResourceMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation CreateResource($input: CreateResourceInput!) {
	resource: createResource(input: $input) {
%s
		turbot: get(path:"turbot")
	}
}`, buildResourceProperties(properties))
}

func updateResourceMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation UpdateResource($input: UpdateResourceInput!) {
 	resource: updateResource(input: $input) {
%s
		turbot: get(path:"turbot")
	}
}`, buildResourceProperties(properties))
}

func putResourceMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation PutResource($input: PutResourceInput!) {
 	resource: putResource(input: $input) {
%s
		turbot: get(path:"turbot")
	}
}`, buildResourceProperties(properties))
}

func deleteResourceMutation() string {
	return `mutation DeleteResource($input: DeleteResourceInput!) {
 	resource: deleteResource(input: $input) {
		turbot: get(path:"turbot")
	}
}`
}

// support properties array of Interface
//
// The id is a GraphQL variable ($id), never interpolated. Only the property selection - built
// internally by buildResourceProperties from a fixed list, never from caller input - is
// interpolated. See TestReadBuildersDoNotInterpolateIdentifier.
func readResourceQuery(properties []interface{}) string {
	return fmt.Sprintf(`query ReadResource($id: ID!) {
	resource(id: $id) {
		type {
			uri
		}
%s
		turbot: get(path:"turbot")
  	}
}`, buildResourceProperties(properties))
}

func getResourceTypeIdQuery() string {
	return `query GetResourceTypeId($id: ID!) {
	resource(id: $id) {
		turbot {
			resourceTypeId
		}
  	}
}`
}

func readResourceListQuery(filter string, properties map[string]string) string {
	var propertiesString bytes.Buffer
	if properties != nil {
		for alias, propertyPath := range properties {
			propertiesString.WriteString(fmt.Sprintf("\t\t\t%s: get(path: \"%s\")\n", alias, propertyPath))
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

func readFullResourceQuery() string {
	return `query ReadFullResource($id: ID!) {
  resource(id: $id) {
	type {
		uri
	}
    data
    turbot: get(path:"turbot")
  }
}`
}

// google directory read query
func readGoogleDirectoryQuery() string {
	return `query ReadGoogleDirectory($id: ID!) {
	directory: resource(id: $id) {
		title:             	get(path:"title")
		parent:            	get(path:"turbot.parentId")
		description:       	get(path:"description")
		status:            	get(path:"status")
		directoryType:     	get(path:"directoryType")
		profileIdTemplate: 	get(path:"profileIdTemplate")
		clientID:          	get(path:"clientID")
		poolId:            	get(path:"poolId")
		groupIdTemplate:   	get(path:"groupIdTemplate")
		loginNameTemplate: 	get(path:"loginNameTemplate")
		hostedDomain:       get(path:"hostedDomain")
		turbot: 			get(path:"turbot")
	}
}`
}

// grant
// The id is a GraphQL variable; only the metadata fragment (a fixed internal selection) is
// interpolated. See TestReadBuildersDoNotInterpolateIdentifier.
func readGrantQuery() string {
	return fmt.Sprintf(`query ReadGrant($id: ID!) {
	grant: grant(id: $id) {
		permissionTypeId
		permissionLevelId
		%s
	}
  }`, turbotGrantMetadataFragment("\t\t"))
}

func createGrantMutation() string {
	return fmt.Sprintf(`mutation CreateGrant($input: CreateGrantInput!) {
	grants: createGrant(input: $input) {
%s
	}
}`, turbotGrantMetadataFragment("\t\t\t"))
}

func deleteGrantMutation() string {
	return fmt.Sprintf(`mutation DeleteGrant($input: DeleteGrantInput!) {
 	grant: deleteGrant(input: $input) {
%s
	}
}`, turbotGrantMetadataFragment("\t\t\t"))
}

// active grant
func readActiveGrantQuery() string {
	return fmt.Sprintf(`query ReadActiveGrant($id: ID!) {
	activeGrant: activeGrant(id: $id){
%s
	}
}`, turbotActiveGrantMetadataFragment("\t\t"))
}

func activateGrantMutation() string {
	return fmt.Sprintf(`mutation ActivateGrant($input: ActivateGrantInput!) {
	grantActivate: activateGrant(input: $input) {
%s
	}
}`, turbotActiveGrantMetadataFragment("\t\t\t"))
}

func deactivateGrantMutation() string {
	return fmt.Sprintf(`mutation DeactivateGrant($input: DeactivateGrantInput!) {
	deactivateGrant(input: $input) {
%s
	}
}`, turbotActiveGrantMetadataFragment("\t\t\t"))
}

// turbot directory
func createTurbotDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createTurbotDirectory($input: CreateTurbotDirectoryInput!) {
 	 	resource: createTurbotDirectory(input: $input){
%s
    	turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateTurbotDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateTurbotDirectory($input: UpdateTurbotDirectoryInput!) {
  		resource: updateTurbotDirectory(input: $input){
%s
		turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func buildResourceProperties(resourceProperties []interface{}) string {
	var propertiesString bytes.Buffer
	if resourceProperties != nil {
		for _, propertyPath := range resourceProperties {
			property, ok := propertyPath.(map[string]string)
			if ok {
				for alias, property := range property {
					propertiesString.WriteString(fmt.Sprintf("\t\t\t%s: get(path: \"%s\")\n", alias, property))
				}
			} else {
				propertiesString.WriteString(fmt.Sprintf("\t\t\t%s: get(path: \"%s\")\n", propertyPath, propertyPath))
			}

		}
	}
	return propertiesString.String()
}

// local directory
func createLocalDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createLocalDirectory($input: CreateLocalDirectoryInput!) {
 	 	resource: createLocalDirectory(input: $input){
%s
    	turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateLocalDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateLocalDirectory($input: UpdateLocalDirectoryInput!) {
  		resource: updateLocalDirectory(input: $input){
%s
		turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

// google directory
func createGoogleDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createGoogleDirectory($input: CreateGoogleDirectoryInput!) {
  		resource: createGoogleDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateGoogleDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateGoogleDirectory($input: UpdateGoogleDirectoryInput!) {
  		resource: updateGoogleDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

// saml directory
func createSamlDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createSamlDirectory($input: CreateSamlDirectoryInput!) {
  		resource: createSamlDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateSamlDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateSamlDirectory($input: UpdateSamlDirectoryInput!) {
  		resource: updateSamlDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

// control
func readControlQuery(args string) string {
	return fmt.Sprintf(`{
control(%s){
	type{
		uri
	}
	state
	mute
	reason
	details
	turbot {
		id
		resourceId
	}
}
}`, args)
}

func muteControlMutation() string {
	return `mutation MuteControl($input: MuteControlInput!) {
		muteControl: muteControl(input: $input) {
			mute
			type {
				uri
			}
			state
			turbot {
				id
				resourceId
			}
		}
	}`
}

func unMuteControlMutation() string {
	return `mutation UnmuteControl($input: UnmuteControlInput!) {
		unmuteControl: unmuteControl(input: $input) {
			mute
			type {
				uri
			}	
			state
			turbot {
				id
				resourceId
			}
		}
	}`
}

// group profile
func createGroupProfileMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createGroupProfile($input: CreateGroupProfileInput!) {
  		resource: createGroupProfile(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateGroupProfileMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateGroupProfile($input: UpdateGroupProfileInput!) {
  		resource: updateGroupProfile(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func deleteGroupProfileMutation() string {
	return fmt.Sprintf(`mutation deleteGroupProfile($input: DeleteGroupProfileInput!) {
  		resource: deleteGroupProfile(input: $input){
    	turbot:get(path:"turbot")
  }
}`)
}

// ldap directory
func createLdapDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createLdapDirectory($input: CreateLdapDirectoryInput!) {
  		resource: createLdapDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateLdapDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateLdapDirectory($input: UpdateLdapDirectoryInput!) {
  		resource: updateLdapDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func deleteLdapDirectory() string {
	return `mutation DeleteResource($input: DeleteResourceInput!) {
 	resource: deleteResource(input: $input) {
		turbot: get(path:"turbot")
	}
}`
}

// get turbot workspace version
func (client *Client) GetTurbotWorkspaceVersion() (*semver.Version, error) {
	query := readPolicyValueQuery("tmod:@turbot/turbot#/policy/types/workspaceVersion", "tmod:@turbot/turbot#/")
	responseData := &PolicyValueResponse{}
	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy value: %s", err.Error())
	}
	// convert interface {} to string
	versionValue := fmt.Sprintf("%v", responseData.PolicyValue.Value)
	// convert version value to semver value
	version, err := semver.New(versionValue)
	if err != nil {
		return nil, fmt.Errorf("error reading guardrails workspace version value: %s", err.Error())
	}
	return version, nil
}

// readPolicyPackIdentityQuery resolves a policy pack by id or aka.
//
// Uses `policyPack(id:)` rather than `resource(id:)` on purpose: policy packs live at the
// Turbot root by default, and `resource(id:)` on a pack requires a grant wherever it sits.
// `policyPack(id:)` is the query the Guardrails console uses and is authorized for an identity
// holding permissions only on the attachment target. It accepts either a numeric id or an aka.
//
// The identifier is passed as a GraphQL variable rather than interpolated, matching what every
// mutation in this client already does. Interpolating it would let a config-supplied value
// containing a double quote escape the string literal and append arbitrary GraphQL, executed with
// the provider's credentials.
func readPolicyPackIdentityQuery() string {
	return `query ReadPolicyPackIdentity($id: ID!) {
	policyPack(id: $id) {
		turbot {
			id
			akas
		}
	}
}`
}

// readAttachedPolicyPacksQuery lists the policy packs attached to a resource, read from the
// resource side so it only needs permissions on that resource. Accepts a numeric id or an aka,
// passed as a GraphQL variable - see readPolicyPackIdentityQuery.
func readAttachedPolicyPacksQuery() string {
	// notFound: RETURN_NULL makes a missing target arrive as `resource: null` instead of an error,
	// while "all other errors are returned" per the schema. That lets Exists distinguish "the target
	// is gone" from "the read failed" by response SHAPE, rather than by matching error text - see
	// ErrTargetNotFound.
	return `query ReadAttachedPolicyPacks($id: ID!) {
	resource(id: $id, options: {notFound: RETURN_NULL}) {
		attachedSmartFolders {
			paging {
				next
			}
			items {
				turbot {
					id
					akas
				}
			}
		}
	}
}`
}
