package apiclient

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
)

func (client *Client) CreateResource(typeAka, parentAka, payload string) (*TurbotMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}

	commandPayload := map[string]string{}
	json.Unmarshal([]byte(payload), &commandPayload)

	commandMeta := map[string]string{
		"typeAka":   typeAka,
		"parentAka": parentAka,
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

func (client *Client) ReadResource(id string, properties map[string]string) (*Resource, error) {
	log.Println("[INFO] ReadResource", id)
	query := readResourceQuery(id, properties)
	var responseData = &ReadResourceResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading resource: %s", err.Error())
	}

	resource, err := client.AssignResourceResults(responseData.Resource, properties)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (client *Client) DeleteResource(aka string) error {
	query := deleteResourceMutation()
	var responseData interface{}

	commandPayload := map[string]string{
		"aka": aka,
	}
	commandMeta := map[string]string{}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting folder: %s", err.Error())
	}
	return nil
}

func (client *Client) UpdateResource(id, parent, title, description string) (*TurbotMetadata, error) {
	query := updateResourceMutation()
	responseData := &UpdateResourceResponse{}
	var commandPayload = map[string]map[string]interface{}{
		"data": {
			"title":       title,
			"description": description,
		},
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
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ResourceExists(id string) (bool, error) {
	resource, err := client.ReadResource(id, nil)
	if err != nil {
		return false, err
	}
	exists := resource.Turbot.Id != ""
	return exists, nil
}

func (client *Client) GetResourceAkas(id string) ([]string, error) {
	resource, err := client.ReadResource(id, nil)
	if err != nil {
		return nil, err
	}
	return resource.Turbot.Akas, nil
}

// assign the ReadResource results into a Resource object, based on the 'properties' map
func (client *Client) AssignResourceResults(responseData interface{}, properties map[string]string) (*Resource, error) {
	var resource Resource
	// initialise map
	resource.Data = make(map[string]interface{})
	// convert turbot property to structure
	if err := mapstructure.Decode(responseData.(map[string]interface{})["turbot"], &resource.Turbot); err != nil {
		return nil, err
	}
	if properties != nil {
		for p := range properties {
			resource.Data[p] = responseData.(map[string]interface{})[p]
		}
	}
	return &resource, nil

}
