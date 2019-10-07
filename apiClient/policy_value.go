package apiClient

import (
	"fmt"
)

func (client *Client) ReadPolicyValue(policyTypeUri, resourceAka string) (*PolicyValue, error) {
	query := readPolicyValueQuery(policyTypeUri, resourceAka)
	responseData := &PolicyValueResponse{}
	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy value: %s", err.Error())
	}

	return &responseData.PolicyValue, nil
}
