package apiclient

import (
	"fmt"
	"log"
)

func (client *Client) CreateSmartFolder(parent string, data map[string]interface{}) (*TurbotMetadata, error) {
	query := createSmartFolderMutation()
	responseData := &CreateSmartFolderResponse{}
	var commandPayload = map[string]interface{}{
		"data": data,
	}
	commandMeta := map[string]string{
		"parentAka": parent,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}

	log.Println("resource", variables)

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.SmartFolder.Turbot, nil
}

func (client *Client) ReadSmartFolder(id string) (*SmartFolder, error) {
	query := readSmartFolderQuery(id)
	//var responseData interface{}
	responseData := &ReadSmartFolderResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.SmartFolder, nil
}

func (client *Client) UpdateSmartFolder(id, parent string, data map[string]interface{}) (*TurbotMetadata, error) {
	query := updateSmartFolderMutation()
	responseData := &UpdateSmartFolderResponse{}
	//var responseData interface{}
	var commandPayload = map[string]map[string]interface{}{
		"data": data,
		"turbotData": {
			"akas": []string{id},
		},
	}
	commandMeta := map[string]interface{}{
		"smartFolderId": id,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating smart folder: %s", err.Error())
	}
	return &responseData.SmartFolder.Turbot, nil
}
