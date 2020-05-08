package apiClient

import (
	"fmt"
)

// create a map of the properties we want the graphql query to return
var localDirectoryUserProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"email",
	"status",
	"displayName",
	"givenName",
	"middleName",
	"familyName",
	"picture",
}

func (client *Client) CreateLocalDirectoryUser(input map[string]interface{}) (*LocalDirectoryUser, error) {
	query := createResourceMutation(localDirectoryUserProperties)
	responseData := &LocalDirectoryUserResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot-iam#/resource/types/localDirectoryUser"
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating local directory user: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadLocalDirectoryUser(id string) (*LocalDirectoryUser, error) {

	query := readResourceQuery(id, localDirectoryUserProperties)
	responseData := &LocalDirectoryUserResponse{}
	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading local directory user: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateLocalDirectoryUserResource(input map[string]interface{}) (*LocalDirectoryUser, error) {
	query := updateResourceMutation(localDirectoryUserProperties)
	responseData := &LocalDirectoryUserResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating local directory user: %s", err.Error())
	}
	return &responseData.Resource, nil
}
