package apiClient

import (
	"fmt"
)

func (client *Client) CreateSmartFolder(input map[string]interface{}) (*SmartFolder, error) {
	query := createSmartFolderMutation()
	responseData := &SmartFolderResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating smart folder: %s", err.Error())
	}
	return &responseData.SmartFolder, nil
}

func (client *Client) ReadSmartFolder(id string) (*SmartFolder, error) {
	query := readSmartFolderQuery(id)
	responseData := &SmartFolderResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading smart folder: %s", err.Error())
	}
	return &responseData.SmartFolder, nil
}

func (client *Client) UpdateSmartFolder(input map[string]interface{}) (*SmartFolder, error) {
	query := updateSmartFolderMutation()
	responseData := &SmartFolderResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating smart folder: %s", err.Error())
	}
	return &responseData.SmartFolder, nil
}
