package apiClient

import (
	"fmt"
)

func (client *Client) CreateProfile(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot-iam#/resource/types/profile"
	variables := map[string]interface{}{
		"input": input,
	}
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

func (client *Client) UpdateProfile(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation()
	responseData := &UpdateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating profile: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}
