package apiClient

import (
	"fmt"
)

func (client *Client) CreateSmartFolder(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createSmartFolderMutation()
	responseData := &CreateSmartFolderResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

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

func (client *Client) UpdateSmartFolder(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateSmartFolderMutation()
	responseData := &UpdateSmartFolderResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating smart folder: %s", err.Error())
	}
	return &responseData.SmartFolder.Turbot, nil
}
