package apiclient

import (
	"fmt"
)

func (client *Client) CreateSmartFolderAttachment(resourceId string, resourceGroupId string) (*TurbotMetadata, error) {
	query := createSmartFolderAttachmentMutation()
	responseData := &CreateResourceResponse{}
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
	return &responseData.Resource.Turbot, nil
}

func (client *Client) ReadSmartFolderAttachment(id string) (*SmartFolderAttachment, error) {
	// create a map of the properties we want the graphql query to return
	properties := map[string]string{
		"resource_id":       "resource_group_id",
		"resource_group_id": "resource_group",
	}
	query := readResourceQuery(id, properties)
	responseData := &ReadSmartFolderAttachmentResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}
