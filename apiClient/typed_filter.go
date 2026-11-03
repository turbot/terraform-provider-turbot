package apiClient

import (
	"errors"
	"fmt"
)

var ErrTypedFilterNotFound = errors.New("typed filter not found")

func IsTypedFilterNotFound(err error) bool {
	return errors.Is(err, ErrTypedFilterNotFound)
}

func (client *Client) CreateTypedFilter(input map[string]interface{}) (*TypedFilter, error) {
	responseData := &TypedFilterResponse{}
	variables := map[string]interface{}{"input": input}

	if err := client.doRequest(createTypedFilterMutation(), variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "typed filter")
	}
	if responseData.TypedFilter == nil {
		return nil, errors.New("error creating typed filter: the server returned no filter")
	}
	return responseData.TypedFilter, nil
}

// ReadTypedFilter returns ErrTypedFilterNotFound when the id names nothing, or names a resource that
// is not a typed filter (the API answers such an id with a null type rather than an error).
func (client *Client) ReadTypedFilter(id string) (*TypedFilter, error) {
	responseData := &TypedFilterResponse{}

	if err := client.doRequest(readTypedFilterQuery(), map[string]interface{}{"id": id}, responseData); err != nil {
		return nil, fmt.Errorf("error reading typed filter: %s", err.Error())
	}
	if responseData.TypedFilter == nil || responseData.TypedFilter.Type == "" {
		return nil, fmt.Errorf("%w: %s", ErrTypedFilterNotFound, id)
	}
	return responseData.TypedFilter, nil
}

func (client *Client) UpdateTypedFilter(input map[string]interface{}) error {
	var responseData interface{}
	variables := map[string]interface{}{"input": input}

	if err := client.doRequest(updateTypedFilterMutation(), variables, &responseData); err != nil {
		return client.handleUpdateError(err, input, "typed filter")
	}
	return nil
}

// DeleteTypedFilter reports an already-deleted filter as success. The API answers a delete of a
// missing id with BAD_USER_INPUT rather than NOT_FOUND, so a failed delete is followed by a read to
// tell "already gone" from a real refusal (for example a filter still used by a guardrail).
func (client *Client) DeleteTypedFilter(id string) error {
	var responseData interface{}

	deleteErr := client.doRequest(deleteTypedFilterMutation(), map[string]interface{}{"id": id}, &responseData)
	if deleteErr == nil {
		return nil
	}
	if _, readErr := client.ReadTypedFilter(id); IsTypedFilterNotFound(readErr) {
		return nil
	}
	return fmt.Errorf("error deleting typed filter: %s", deleteErr.Error())
}
