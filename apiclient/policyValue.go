package apiclient

import (
	"errors"
	"fmt"
	"log"
)

func (client *Client) ReadPolicyValue(policyTypeUri, resourceAka string) (*PolicyValue, error) {
	query := readPolicyValueQuery(policyTypeUri, resourceAka)
	responseData := &PolicyValueResponse{}

	log.Println("ReadPolicyValue")
	log.Println("policy_type", policyTypeUri)
	log.Println("resource", resourceAka)
	log.Println("Query:", query)

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy value: %s", err.Error())
	}

	log.Println("DoRequest returned:", responseData)
	return &responseData.PolicyValue, nil
}

func validateGetPolicyValueArgs(policyTypeUri, resourceAka, policyValueId string) error {
	if policyValueId != "" && (policyTypeUri != "" || resourceAka != "") ||
		(policyValueId == "" && (policyTypeUri == "" || resourceAka == "")) {
		return errors.New("GetPolicyValueById must be called with either policyValueId, or policyTypeUri AND resourceAka")
	}
	return nil
}
