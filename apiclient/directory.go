package apiclient

import (
	"fmt"
)

func (client *Client) CreateLocalDirectory(parent, title, description string, profileId string, status string, directoryType string) (*TurbotMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}
	var commandPayload = map[string]map[string]string{
		"data": {
			"title":             title,
			"description":       description,
			"profileIdTemplate": profileId,
			"status":            status,
			"directoryType":     directoryType,
		},
	}
	commandMeta := map[string]string{
		"typeAka":   "tmod:@turbot/turbot-iam#/resource/types/localDirectory",
		"parentAka": parent,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadLocalDirectory(id string) (*LocalDirectory, error) {
	// create a map of the properties we want the graphql query to return
	properties := map[string]string{
		"title":             "title",
		"parent":            "turbot.parentId",
		"description":       "description",
		"status":            "status",
		"directoryType":     "directoryType",
		"profileIdTemplate": "profileIdTemplate",
	}
	query := readResourceQuery(id, properties)
	responseData := &ReadLocalDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}

// func (client *Client) FindDirectory(title, parentAka string) ([]Directory, error) {
// 	responseData := &FindDirectoryResponse{}
// 	// convert parentAka into an id
// 	parent, err := client.ReadResource(parentAka, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	parentId := parent.Turbot.Id
// 	query := findDirectoryQuery(title, parentId)

// 	// execute api call
// 	if err := client.doRequest(query, nil, &responseData); err != nil {
// 		return nil, fmt.Errorf("error reading Directory: %s", err.Error())
// 	}

// 	return responseData.Directories.Items, nil
// }
