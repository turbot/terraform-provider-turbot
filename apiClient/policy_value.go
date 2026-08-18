package apiClient

func (client *Client) ReadPolicyValue(policyTypeUri, resourceAka string) (*PolicyValue, error) {
	query := readPolicyValueQuery()
	responseData := &PolicyValueResponse{}
	variables := map[string]interface{}{"uri": policyTypeUri, "resourceId": resourceAka}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleReadError(err, policyTypeUri, "policy setting")
	}

	return &responseData.PolicyValue, nil
}
