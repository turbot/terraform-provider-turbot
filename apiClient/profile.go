package apiClient

import (
	"fmt"
	"log"
)

func (client *Client) CreateProfile(parent string, data map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}
	var commandPayload = map[string]interface{}{
		"data": data,
	}
	commandMeta := map[string]string{
		"typeAka":   "tmod:@turbot/turbot-iam#/resource/types/profile",
		"parentAka": parent,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}
	log.Println("[INFO} resourceTurbotProfileCreate", variables)
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating profile: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadProfile(id string) (*Profile, error) {
	// create a map of the properties we want the graphql query to return
	properties := map[string]string{
		"title":           "title",
		"parent":          "turbot.parentId",
		"status":          "status",
		"displayName":     "displayName",
		"email":           "email",
		"givenName":       "givenName",
		"familyName":      "familyName",
		"directoryPoolId": "directoryPoolId",
	}
	query := readResourceQuery(id, properties)
	responseData := &ReadProfileResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading profile: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateProfile(id, parent string, data map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation()
	responseData := &UpdateResourceResponse{}
	var commandPayload = map[string]map[string]interface{}{
		"data": data,
		"turbotData": {
			"akas": []string{id},
		},
	}
	commandMeta := map[string]interface{}{
		"typeAka":   "tmod:@turbot/turbot#/resource/types/folder",
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
		return nil, fmt.Errorf("error creating profile: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}
