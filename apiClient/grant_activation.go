package apiClient

import (
	"fmt"
)

func (client *Client) CreateGrantActivation(grant, resourceAka string) (*TurbotActiveGrantMetadata, error) {
	query := activateGrantMutation()
	responseData := &ActivateGrantResponse{}
	//var responseData interface{}
	commandMeta := map[string]interface{}{
		"resourceAka": resourceAka,
		"grantId":     grant,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"commands": []map[string]interface{}{{
				"meta": commandMeta,
			}},
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating grant: %s", err.Error())
	}
	return &responseData.GrantActivate.Items[0].Turbot, nil
}

func (client *Client) ReadGrantActivation(id string) (*ActiveGrant, error) {
	query := readActiveGrantQuery(id)
	responseData := &ReadActiveGrantResponse{}
	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading folder: %s", err.Error())
	}
	return &responseData.ActiveGrant, nil
}

func (client *Client) DeleteGrantActivation(id string) error {
	query := deactivateGrantMutation()
	var responseData interface{}

	commandMeta := map[string]interface{}{
		"activationId": id,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"commands": []map[string]interface{}{{
				"meta": commandMeta,
			}},
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting folder: %s", err.Error())
	}
	return nil
}

func (client *Client) GrantActivationExists(id string) (bool, error) {
	grantActivate, err := client.ReadGrantActivation(id)
	if err != nil {
		return false, err
	}
	exists := grantActivate.Turbot.Id != ""
	return exists, nil
}
