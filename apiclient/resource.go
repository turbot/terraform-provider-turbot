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

	log.Println("CreateResource")
	log.Println("parentAka", parentAka)
	log.Println("commandPayload", commandPayload)
	log.Println("query:", query)

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	log.Println("DoRequest returned:", responseData)
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadResource(id string, properties map[string]string) (*Resource, error) {
	query := readResourceQuery(id, properties)
	var responseData = &ReadResourceResponse{}

	log.Println("ReadFolder")
	log.Println("id", id)
	log.Print("Query:", query)

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy: %s", err.Error())
	}

	log.Println("DoRequest returned:", responseData)
	resource := client.AssignResourceResults(responseData.Resource, properties)
	return &resource, nil
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

	log.Println("DeleteResource")
	log.Println("Resource aka", aka)
	log.Println("Query:", query)
	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting folder: %s", err.Error())
	}
	log.Println("DoRequest returned:", responseData)
	return nil
}

func (client *Client) GetResourceAkas(id string) ([]string, error) {
	resource, err := client.ReadResource(id, nil)
	if err != nil {
		return nil, err
	}
	return resource.Turbot.Akas, nil
}

// assign the ReadResource results into a Resource object, based on the 'properties' map
func (client *Client) AssignResourceResults(responseData interface{}, properties map[string]string) Resource {
	var resource Resource
	// convert turbot property to structure

	mapstructure.Decode(responseData.(map[string]interface{})["turbot"], &resource.Turbot)
	if properties != nil {
		for p := range properties {
			resource.Data[p] = responseData.(map[string]interface{})[p]
		}
	}
	return resource

}
