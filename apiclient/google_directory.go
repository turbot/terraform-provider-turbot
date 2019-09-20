package apiclient

import (
	"fmt"
)

func (client *Client) CreateGoogleDirectory(parent string, commandPayload map[string]map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}
	commandMeta := map[string]string{
		"typeAka":   "tmod:@turbot/turbot-iam#/resource/types/googleDirectory",
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

func (client *Client) ReadGoogleDirectory(id string) (*GoogleDirectory, error) {
	// create a map of the properties we want the graphql query to return
	properties := map[string]string{
		"title":             "title",
		"parent":            "turbot.parentId",
		"description":       "description",
		"status":            "status",
		"directoryType":     "directoryType",
		"profileIdTemplate": "profileIdTemplate",
		"clientID":          "clientID",
	//	"clientSecret":      "clientSecret",
		"poolId":            "poolId",
		"groupIdTemplate":   "groupIdTemplate",
		"loginNameTemplate": "loginNameTemplate",
		"hostedName":        "hostedName",
	}
	query := readResourceQuery(id, properties)
	responseData := &ReadGoogleDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateGoogleDirectory(id, parent string, commandPayload map[string]map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation()
	responseData := &UpdateResourceResponse{}
	// add akas to turbotData
	commandPayload["turbotData"]["akas"] = []string{id}
	commandMeta := map[string]interface{}{
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
