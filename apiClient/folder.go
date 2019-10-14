package apiClient

import (
	"fmt"
)

func (client *Client) CreateFolder(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot#/resource/types/folder"
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadFolder(id string) (*Folder, error) {
	// create a map of the properties we want the graphql query to return
	properties := map[string]string{
		"title":       "title",
		"parent":      "turbot.parentId",
		"description": "description",
	}
	query := readResourceQuery(id, properties)
	responseData := &ReadFolderResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateFolder(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation()
	responseData := &UpdateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}
