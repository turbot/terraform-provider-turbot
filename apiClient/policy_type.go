package apiClient

import "fmt"

func (client *Client) FindPolicyType(policyTypeUri string) (PolicyType, error) {
	responseData := &FindPolicyTypeResponse{}

	query := findPolicyTypeQuery()
	variables := map[string]interface{}{
		"filter": []string{fmt.Sprintf("policyTypeId:%s level:self", policyTypeUri)},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return PolicyType{}, client.handleReadError(err, policyTypeUri, "policy type")
	}

	if len(responseData.PolicyTypes.Items) > 0 {
		return responseData.PolicyTypes.Items[0], nil
	}

	return PolicyType{}, nil
}
