package apiclient

import (
	"fmt"
)

func (client *Client) CreateGrant(profileId, parent string, data map[string]interface{}) (*TurbotGrantMetadata, error) {
	query := createGrantMutation()
	responseData := &CreateGrantResponse{}
	commandMeta := map[string]interface{}{
		"resourceAka": parent,
		"profileId":   profileId,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"commands": []map[string]interface{}{{
				"payload": data,
				"meta":    commandMeta,
			}},
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating grant: %s", err.Error())
	}
	//ret := TurbotGrantMetadata{
	//	Id: "testID",
	//}
	return &responseData.Grants.Items[0].Turbot, nil
}

func (client *Client) ReadGrant(id string) (*Grant, error) {
	query := readGrantQuery(id)
	responseData := &ReadGrantResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Grant, nil
}
