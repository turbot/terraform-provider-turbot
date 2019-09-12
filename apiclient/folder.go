package apiclient

import (
	"fmt"
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

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
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

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) FindFolder(title, parentAka string) ([]Folder, error) {
	responseData := &FindFolderResponse{}
	// convert parentAka into an id
	parent, err := client.ReadResource(parentAka, nil)
	if err != nil {
		return nil, err
	}
	parentId := parent.Turbot.Id
	query := findFolderQuery(title, parentId)

	// execute api call
	if err := client.doRequest(query, nil, &responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}

	return responseData.Folders.Items, nil
}
