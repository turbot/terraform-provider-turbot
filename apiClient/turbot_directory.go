package apiClient

import (
	"fmt"
)

var turbotDirectoryProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId", "tags": "turbot.tags"},
	"title",
	"description",
	"status",
	"directoryType",
	"profileIdTemplate",
	"server",
}

func (client *Client) CreateTurbotDirectory(input map[string]interface{}) (*TurbotDirectory, error) {
	query := createTurbotDirectoryMutation(turbotDirectoryProperties)
	responseData := &TurbotDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating turbot directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadTurbotDirectory(id string) (*TurbotDirectory, error) {
	// create a map of the properties we want the graphql query to return
	query := readResourceQuery(id, turbotDirectoryProperties)
	responseData := &TurbotDirectoryResponse{}
	// execute api call
	if err := client.doRequest(query, nil, &responseData); err != nil {
		return nil, fmt.Errorf("error reading turbot directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateTurbotDirectory(input map[string]interface{}) (*TurbotDirectory, error) {
	query := updateTurbotDirectoryMutation(turbotDirectoryProperties)
	responseData := &TurbotDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return nil, fmt.Errorf("error updating turbot directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}
