package apiClient

func (client *Client) CreateRollout(input map[string]interface{}) (*Rollout, error) {
	query := createRolloutMutation()
	responseData := &RolloutResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "rollout")
	}
	return &responseData.Rollout, nil
}

func (client *Client) ReadRollout(id string) (*Rollout, error) {
	query := readRolloutQuery(id)
	responseData := &RolloutResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "rollout")
	}
	return &responseData.Rollout, nil
}

func (client *Client) UpdateRollout(input map[string]interface{}) (*Rollout, error) {
	query := updateRolloutMutation()
	responseData := &RolloutResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "rollout")
	}
	return &responseData.Rollout, nil
}
