package apiclient

import (
	"fmt"
	"log"
)

func (client *Client) CreatePolicySetting(policyTypeUri, resourceAka string, commandPayload interface{}) (*PolicySetting, error) {
	query := createPolicySettingMutation()
	responseData := &PolicySettingResponse{}

	commandMeta := map[string]string{
		"policyTypeUri": policyTypeUri,
		"resourceAka":   resourceAka,
	}

	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}

	log.Println("CreatePolicySetting")
	log.Println("policy_type", policyTypeUri)
	log.Println("resource", resourceAka)
	log.Println("Query:", query)

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating policy: %s", err.Error())
	}
	log.Println("DoRequest returned:", responseData)
	return &responseData.PolicySetting, nil
}

func (client *Client) ReadPolicySetting(id string) (*PolicySetting, error) {
	query := readPolicySettingQuery(id)
	responseData := &PolicySettingResponse{}

	log.Println("ReadPolicySetting")
	log.Println("id", id)
	log.Print("Query:", query)

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy: %s", err.Error())
	}

	log.Println("DoRequest returned:", responseData)
	return &responseData.PolicySetting, nil
}

func (client *Client) UpdatePolicySetting(id string, commandPayload interface{}) error {
	query := updatePolicySettingMutation()
	responseData := &PolicySettingResponse{}

	commandMeta := map[string]string{
		"policySettingId": id,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"payload": commandPayload,
			"meta":    commandMeta,
		},
	}

	log.Println("UpdatePolicySetting")
	log.Println("id", id)
	log.Println("Query:", query)

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error updating policy: %s", err.Error())
	}
	log.Println("DoRequest returned:", responseData)
	return nil
}

func (client *Client) DeletePolicySetting(id string) error {
	query := deletePolicySettingMutation()
	responseData := &PolicySettingResponse{}

	commandMeta := map[string]string{
		"policySettingId": id,
	}
	variables := map[string]interface{}{
		"command": map[string]interface{}{
			"meta": commandMeta,
		},
	}

	log.Println("DeletePolicySetting")
	log.Println("id", id)
	log.Println("Query:", query)

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting policy: %s", err.Error())
	}
	log.Println("DoRequest returned:", responseData)
	return nil
}
