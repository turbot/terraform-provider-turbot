package apiClient

import (
	"fmt"
)

// create a map of the properties we want the graphql query to return
var samlDirectoryProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"description",
	"status",
	"directoryType",
	"profileIdTemplate",
}

func (client *Client) CreateSamlDirectory(input map[string]interface{}) (*SamlDirectory, error) {
	query := createResourceMutation(samlDirectoryProperties)
	responseData := &SamlDirectoryResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot-iam#/resource/types/samlDirectory"
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadSamlDirectory(id string) (*SamlDirectory, error) {

	query := readResourceQuery(id, samlDirectoryProperties)
	responseData := &SamlDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateSamlDirectory(input map[string]interface{}) (*SamlDirectory, error) {
	query := updateResourceMutation(samlDirectoryProperties)
	responseData := &SamlDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}
