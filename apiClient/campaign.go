package apiClient

func (client *Client) CreateCampaign(input map[string]interface{}) (*Campaign, error) {
	query := createCampaignMutation()
	responseData := &CampaignResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "campaign")
	}
	return &responseData.Campaign, nil
}

func (client *Client) ReadCampaign(id string) (*Campaign, error) {
	query := readCampaignQuery(id)
	responseData := &CampaignResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "campaign")
	}
	return &responseData.Campaign, nil
}

func (client *Client) UpdateCampaign(input map[string]interface{}) (*Campaign, error) {
	query := updateCampaignMutation()
	responseData := &CampaignResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "campaign")
	}
	return &responseData.Campaign, nil
}
