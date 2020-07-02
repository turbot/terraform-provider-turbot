package apiClient

import (
	"fmt"
)

// legacy google directory resource properties
var googleDirectoryPropertiesLegacy = []interface{}{
	// explicit mapping
	map[string]string{"client_id": "clientID"},
	// implicit mappings
	"title", "poolId", "profileIdTemplate", "groupIdTemplate", "loginNameTemplate", "clientSecret", "hostedName", "description"}

var googleDirectoryProperties = []interface{}{
	// explicit mapping
	map[string]string{"tags": "turbot.tags"},
	// implicit mappings
	"title", "poolId", "profileIdTemplate", "groupIdTemplate", "loginNameTemplate", "clientSecret", "hostedDomain", "description", "clientId"}

// legacy create/update functions
func (client *Client) CreateGoogleDirectoryLegacy(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createResourceMutation(googleDirectoryPropertiesLegacy)
	responseData := &CreateResourceResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot-iam#/resource/types/googleDirectory"
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating google directory: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) UpdateGoogleDirectoryLegacy(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation(googleDirectoryPropertiesLegacy)
	responseData := &UpdateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating google directory: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadGoogleDirectory(id string) (*GoogleDirectory, error) {
	/*
		GoogleDirectory read response has clientSecret attribute,
		which is fetched from getSecret(path:"clientSecret") and
		not from get() resolver.
		That's why we used separate query and not readResourceQuery()
	*/
	query := readGoogleDirectoryQuery(id)
	responseData := &ReadGoogleDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading google directory: %s", err.Error())
	}
	return &responseData.Directory, nil
}

func (client *Client) CreateGoogleDirectory(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createGoogleDirectoryMutation(googleDirectoryProperties)
	responseData := &CreateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating google directory: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) UpdateGoogleDirectory(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateGoogleDirectoryMutation(googleDirectoryProperties)
	responseData := &UpdateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating google directory: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}
