package apiclient

import (
	"fmt"
)

func (client *Client) CreateSmartFolderAttachment(resourceId string, resourceGroupId string) (*TurbotMetadata, error) {
	query := createSmartFolderAttachmentMutation()
	responseData := &CreateSmartFolderAttachResponse{}
	commandMeta := map[string]string{
		"resourceId":      resourceId,
		"resourceGroupId": resourceGroupId,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"meta": commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating folder: %s", err.Error())
	}
	return &responseData.SmartFolder.Turbot, nil
}

func (client *Client) ReadSmartFolderAttachment(id string) (*[]Resoures, error) {
	query := readAttachedResourcesOnSmartfolder(id)
	responseData := &ReadSmartFolderAttachResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.SmartFolder.Items, nil
}

func (client *Client) DeleteSmartFolderAttachment(resourceId string, resourceGroupId string) error {
	_ := detachSmartFolderAttachment()
	var responseData interface{}
	commandMeta := map[string]string{
		"resourceId":      resourceId,
		"resourceGroupId": resourceGroupId,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"meta": commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error creating folder: %s", err.Error())
	}
	return nil
}
