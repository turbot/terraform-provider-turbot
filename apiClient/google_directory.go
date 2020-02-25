package apiClient

import (
	"fmt"
)

var googleDirectoryProperties = []interface{}{
	// explicit mapping
	map[string]string{"client_id": "clientID"},
	// implicit mappings
	"title", "poolId", "profileIdTemplate", "groupIdTemplate", "loginNameTemplate", "clientSecret", "hostedName", "description"}

func (client *Client) CreateGoogleDirectory(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createResourceMutation(googleDirectoryProperties)
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

func (client *Client) UpdateGoogleDirectory(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation(googleDirectoryProperties)
	responseData := &UpdateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating google directory: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}
