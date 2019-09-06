package apiclient

import (
	"fmt"
	"log"
)

func (client *Client) CreateFolder(parent, title, description string) (*TurbotMetadata, error) {
	query := createResourceMutation()
	responseData := &CreateResourceResponse{}

	var commandPayload = map[string]map[string]string{
		"data": {
			"title":       title,
			"description": description,
		},
	}
	commandMeta := map[string]string{
		"typeAka":   "tmod:@turbot/turbot#/resource/types/folder",
		"parentAka": parent,
	}

	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}

	log.Println("CreateFolder")
	log.Println("parentAka", parent)
	log.Println("title", title)
	log.Println("description", description)
	log.Println("query:", query)

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	log.Println("DoRequest returned:", responseData)
	// set
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadFolder(id string) (*Folder, error) {
	// create a map of the properties we want the graphql query to return
	properties := map[string]string{
		"title":       "title",
		"parent":      "turbot.parentId",
		"description": "description",
	}

	query := readResourceQuery(id, properties)
	responseData := &ReadFolderResponse{}

	log.Println("ReadFolder")
	log.Println("id", id)
	log.Print("Query:", query)

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy: %s", err.Error())
	}

	log.Println("DoRequest returned:", responseData)
	return &responseData.Resource, nil
}
