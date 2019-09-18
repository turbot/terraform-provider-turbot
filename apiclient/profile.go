package apiclient

import (
	"fmt"
	"log"
)

func (client *Client) CreateProfile(payload *ProfilePayload) (*TurbotMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}
	var commandPayload = map[string]map[string]string{
		"data": {
			"title":           payload.Title,
			"status":          payload.Status,
			"displayName":     payload.DisplayName,
			"email":           payload.Email,
			"givenName":       payload.GivenName,
			"familyName":      payload.FamilyName,
			"directoryPoolId": payload.DirectoryPoolId,
			"profileId":       payload.ProfileId,
		},
	}
	commandMeta := map[string]string{
		"typeAka":   "tmod:@turbot/turbot-iam#/resource/types/profile",
		"parentAka": payload.Parent,
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

func (client *Client) UpdateProfile(id string, payload *ProfilePayload) (*TurbotMetadata, error) {
	query := updateResourceMutation()
	responseData := &UpdateResourceResponse{}
	var commandPayload = map[string]map[string]interface{}{
		"data": {
			"title":           payload.Title,
			"status":          payload.Status,
			"displayName":     payload.DisplayName,
			"email":           payload.Email,
			"givenName":       payload.GivenName,
			"familyName":      payload.FamilyName,
			"directoryPoolId": payload.DirectoryPoolId,
		},
		"turbotData": {
			"akas": []string{id},
		},
	}
	commandMeta := map[string]interface{}{
		"typeAka":   "tmod:@turbot/turbot#/resource/types/folder",
		"parentAka": payload.Parent,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}
	log.Println("[INFO} resourceTurbotProfileUpdate", variables)
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}
