package apiclient

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
)

func (client *Client) CreateResource(typeAka, parentAka, payload string) (*TurbotResourceMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}

	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return nil, fmt.Errorf("error creating resource: %s", err.Error())
	}

	// todo extract turbotData

	commandMeta := map[string]string{
		"typeAka":   typeAka,
		"parentAka": parentAka,
	}

	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": map[string]interface{}{
				"data": data,
			},
			"meta": commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating resource: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadResource(aka string, properties map[string]string) (*Resource, error) {
	query := readResourceQuery(aka, properties)
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

// todo replace with empty get()
func (client *Client) ReadFullResource(aka string) (*FullResource, error) {
	query := readFullResourceQuery(aka)
	var responseData = &ReadFullResourceResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading resource: %s", err.Error())
	}

	return &responseData.Resource, nil
}

func (client *Client) UpdateResource(id, typeAka, parentAka, payload string) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation()
	responseData := &UpdateResourceResponse{}
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return nil, fmt.Errorf("error updating resource: %s", err.Error())
	}

	// todo extract turbotData

	commandMeta := map[string]string{
		"typeAka":   typeAka,
		"parentAka": parentAka,
	}

	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": map[string]interface{}{
				"data": data,
			},
			"meta": commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.Resource.Turbot, nil
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

func (client *Client) ResourceExists(id string) (bool, error) {
	resource, err := client.ReadResource(id, nil)
	if err != nil {
		return false, err
	}
	exists := resource.Turbot.Id != ""
	return exists, nil
}

func (client *Client) GetResourceAkas(ResourceId string) ([]string, error) {
	resource, err := client.ReadResource(ResourceId, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to load target resource; %s", err)
		return nil, err
	}
	resource_Akas := resource.Turbot.Akas
	// if this resource has no akas, just use the id
	if resource_Akas == nil {
		resource_Akas = []string{ResourceId}
	}
	return resource_Akas, nil
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
	// write properties into a map
	if properties != nil {
		for p := range properties {
			resource.Data[p] = responseData.(map[string]interface{})[p]
		}
	}
	return &resource, nil

}
