package apiClient

import (
	"errors"
	"fmt"
)

var ErrGuardrailNotFound = errors.New("guardrail not found")

func IsGuardrailNotFound(err error) bool {
	return errors.Is(err, ErrGuardrailNotFound)
}

func (client *Client) CreateGuardrail(input map[string]interface{}) (string, error) {
	responseData := &GuardrailResponse{}
	variables := map[string]interface{}{"input": input}

	if err := client.doRequest(createGuardrailMutation(), variables, responseData); err != nil {
		return "", fmt.Errorf("error creating guardrail: %s", err.Error())
	}
	if responseData.Guardrail == nil || responseData.Guardrail.Turbot.Id == "" {
		return "", errors.New("error creating guardrail: the server returned no guardrail")
	}
	return responseData.Guardrail.Turbot.Id, nil
}

// ReadGuardrail returns ErrGuardrailNotFound when the id names nothing, or names a resource that is
// not a guardrail.
func (client *Client) ReadGuardrail(id string) (*Guardrail, error) {
	responseData := &GuardrailResponse{}

	if err := client.doRequest(readGuardrailQuery(), map[string]interface{}{"id": id}, responseData); err != nil {
		return nil, fmt.Errorf("error reading guardrail: %s", err.Error())
	}
	guardrail := responseData.Guardrail
	if guardrail == nil || guardrail.Type == nil || guardrail.Type.Uri == "" {
		return nil, fmt.Errorf("%w: %s", ErrGuardrailNotFound, id)
	}
	return guardrail, nil
}

func (client *Client) ReadGuardrailType(uri string) (*GuardrailType, error) {
	responseData := &GuardrailTypeResponse{}

	if err := client.doRequest(readGuardrailTypeQuery(), map[string]interface{}{"id": uri}, responseData); err != nil {
		return nil, fmt.Errorf("error reading guardrail type %s: %s", uri, err.Error())
	}
	if responseData.GuardrailType == nil || responseData.GuardrailType.Uri == "" {
		return nil, fmt.Errorf("guardrail type %s not found. Is the mod installed?", uri)
	}
	return responseData.GuardrailType, nil
}

func (client *Client) UpdateGuardrail(input map[string]interface{}) error {
	var responseData interface{}
	variables := map[string]interface{}{"input": input}

	if err := client.doRequest(updateGuardrailMutation(), variables, &responseData); err != nil {
		return fmt.Errorf("error updating guardrail: %s", err.Error())
	}
	return nil
}

// DeleteGuardrail reports an already-deleted guardrail as success. A delete of a missing id is
// answered with BAD_USER_INPUT rather than NOT_FOUND, so a failed delete is followed by a read.
func (client *Client) DeleteGuardrail(id string) error {
	var responseData interface{}

	deleteErr := client.doRequest(deleteGuardrailMutation(), map[string]interface{}{"id": id}, &responseData)
	if deleteErr == nil {
		return nil
	}
	if _, readErr := client.ReadGuardrail(id); IsGuardrailNotFound(readErr) {
		return nil
	}
	return fmt.Errorf("error deleting guardrail: %s", deleteErr.Error())
}
