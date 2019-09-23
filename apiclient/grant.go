package apiclient

import (
	"fmt"
)

func (client *Client) CreateGrant(profileId, parent string, data map[string]interface{}) (*TurbotGrantMetadata, error) {
	query := createGrantMutation()
	responseData := &CreateGrantResponse{}
	commandMeta := map[string]interface{}{
		"data": data,
	}
	variables := map[string]interface{}{
		"commands": map[string]interface{}{
			"payload": data,
			"meta":    commandMeta,
		},
	}
	// log.Println("variables", variables)
	// d, err := json.Marshal(variables)
	// if err != nil {
	// 	log.Println("json", d)
	// } else {
	// 	log.Println("err", err)
	// 	log.Println("json", string(d))
	// }

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating grant: %s", err.Error())
	}
	return &responseData.Grant.Items[0].Turbot, nil
}

func (client *Client) ReadGrant(id string) (*Grant, error) {

	query := readGrantQuery(id)
	responseData := &ReadGrantResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.Resource, nil
}
