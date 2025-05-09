package apiClient

import (
	"fmt"
)

func (client *Client) AttachGuardrail(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := attachGuardrailMutation()
	responseData := &AttachGuardrailResponse{}

	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "guardrail attachment")
	}
	return &responseData.Turbot, nil
}

func (client *Client) DetachGuardrail(input map[string]interface{}) error {
	query := detachGuardrailMutation()
	var responseData interface{}

	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting guardrail attachment: %s", err.Error())
	}
	return nil
}
