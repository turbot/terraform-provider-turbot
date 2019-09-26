package apiclient

import (
	"fmt"
)

func (client *Client) CreateSmartFolderAttachment(resourceId, resourceGroupId string) (*TurbotMetadata, error) {
	query := createSmartFolderAttachmentMutation()
	responseData := &CreateSmartFolderAttachResponse{}
	//var res interface{}
	commandMeta := map[string]string{
		"resourceId":    resourceId,
		"smartFolderId": resourceGroupId,
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
	return &responseData.SmartFolderAttach.Turbot, nil
}

func (client *Client) DeleteSmartFolderAttachment(resource, smartFolder string) error {
	query := detachSmartFolderAttachment()
	var responseData interface{}
	commandMeta := map[string]string{
		"resourceId":    resource,
		"smartFolderId": smartFolder,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"meta": commandMeta,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting smart folder: %s", err.Error())
	}
	return nil
}
