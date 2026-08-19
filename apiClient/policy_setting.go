package apiClient

import (
	"fmt"

	errorsHandler "github.com/turbot/terraform-provider-turbot/errors"
)

func (client *Client) CreatePolicySetting(input map[string]interface{}) (*PolicySetting, error) {
	query := createPolicySettingMutation()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "policy setting")
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) ReadPolicySetting(id string) (*PolicySetting, error) {
	query := readPolicySettingQuery()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{"id": id}

	// execute api call
	err := client.doRequest(query, variables, responseData)
	if err == nil {
		return &responseData.PolicySetting, nil
	}
	if !errorsHandler.ForbiddenError(err) {
		return nil, client.handleReadError(err, id, "policy setting")
	}

	// A Forbidden here does not necessarily mean the identity may not read the setting: the
	// primary query requests secretValue/secretValueSource, and Guardrails guards those fields
	// with a Turbot/Admin check at the setting's target resource even when the policy type holds
	// no secret. Retry without the secret fields — Turbot/Metadata or above can read everything
	// else. If the fallback fails too, the denial is on the setting itself (or the workspace
	// predates the fallback's type fields), so surface the original error unchanged.
	fallback := &PolicySettingResponse{}
	if fallbackErr := client.doRequest(readPolicySettingWithoutSecretsQuery(), variables, fallback); fallbackErr != nil {
		return nil, client.handleReadError(err, id, "policy setting")
	}
	setting := fallback.PolicySetting
	// For a genuinely secret policy type the plain fields do not carry the real value, and
	// storing them would corrupt state and produce phantom diffs. Refuse with an actionable
	// error instead.
	if setting.Type.Secret || setting.Type.SecretLevel == "CONFIDENTIAL" || setting.Type.SecretLevel == "SECRET" {
		return nil, fmt.Errorf("error reading policy setting: policy type %s is secret - reading its value requires Turbot/Admin on resource %s (Turbot/Owner at the Turbot root for SECRET-level policy types)", setting.Type.Uri, setting.Turbot.ResourceId)
	}
	return &setting, nil
}

func (client *Client) UpdatePolicySetting(input map[string]interface{}) (*PolicySetting, error) {
	query := updatePolicySettingMutation()
	responseData := &PolicySettingResponse{}

	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "policy setting")
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) DeletePolicySetting(id string) error {
	query := deletePolicySettingMutation()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{
		"input": map[string]string{
			"id": id,
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

	query := findPolicySettingQuery()
	variables := map[string]interface{}{
		"filter": []string{fmt.Sprintf("policyType:%s resource:%s", policyTypeUri, resourceAka)},
	}

	// execute api call
	err := client.doRequest(query, variables, &responseData)
	if err != nil && errorsHandler.ForbiddenError(err) {
		// The Forbidden may be the per-field secretValue guard firing on a matched item rather
		// than a denial of the list itself — see ReadPolicySetting. Retry without the secret
		// fields; if that fails too, surface the original error.
		fallbackData := &FindPolicySettingResponse{}
		if fallbackErr := client.doRequest(findPolicySettingWithoutSecretsQuery(), variables, &fallbackData); fallbackErr == nil {
			responseData = fallbackData
			err = nil
		}
	}
	if err != nil {
		return PolicySetting{}, client.handleReadError(err, policyTypeUri, "policy setting")
	}

	for _, setting := range responseData.PolicySettings.Items {
		if setting.Default {
			return setting, nil
		}
	}
	return PolicySetting{}, nil
}
