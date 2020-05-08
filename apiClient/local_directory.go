package apiClient

import (
	"fmt"
)

var localDirectoryProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"description",
	"status",
	"directoryType",
	"profileIdTemplate",
}

func (client *Client) CreateLocalDirectoryLegacy(input map[string]interface{}) (*LocalDirectory, error) {
	query := createResourceMutation(localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot-iam#/resource/types/localDirectory"
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadLocalDirectory(id string) (*LocalDirectory, error) {
	// create a map of the properties we want the graphql query to return
	query := readResourceQuery(id, localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateLocalDirectoryLegacy(input map[string]interface{}) (*LocalDirectory, error) {
	query := updateResourceMutation(localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) CreateLocalDirectory(input map[string]interface{}) (*LocalDirectory, error) {
	query := createLocalDirectoryMutation(localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateLocalDirectory(input map[string]interface{}) (*LocalDirectory, error) {
	query := updateLocalDirectoryMutation(localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}
