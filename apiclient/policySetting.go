package apiclient

import (
	"fmt"
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

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating policy: %s", err.Error())
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) ReadPolicySetting(id string) (*PolicySetting, error) {
	query := readPolicySettingQuery(id)
	responseData := &PolicySettingResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy setting: %s", err.Error())
	}
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
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error updating policy: %s", err.Error())
	}
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
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting policy: %s", err.Error())
	}
	return nil
}

func (client *Client) FindPolicySetting(policyTypeUri, resourceAka string) (PolicySetting, error) {
	responseData := &FindPolicySettingResponse{}

	query := findPolicySettingQuery(policyTypeUri, resourceAka)

	// execute api call
	if err := client.doRequest(query, nil, &responseData); err != nil {
		return PolicySetting{}, fmt.Errorf("error reading folder: %s", err.Error())
	}

	for _, setting := range responseData.PolicySettings.Items {
		if setting.Default {
			return setting, nil
		}
	}
	return PolicySetting{}, nil
}
