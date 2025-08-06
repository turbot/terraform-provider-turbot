package apiClient

func (client *Client) CreateGuardrail(input map[string]interface{}) (*Guardrail, error) {
	query := createGuardrailMutation()
	responseData := &GuardrailResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "guardrail")
	}
	return &responseData.Guardrail, nil
}

func (client *Client) ReadGuardrail(id string) (*Guardrail, error) {
	query := readGuardrailQuery(id)
	responseData := &GuardrailResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "guardrail")
	}
	return &responseData.Guardrail, nil
}

func (client *Client) UpdateGuardrail(input map[string]interface{}) (*Guardrail, error) {
	query := updateGuardrailMutation()
	responseData := &GuardrailResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "guardrail")
	}
	return &responseData.Guardrail, nil
}
