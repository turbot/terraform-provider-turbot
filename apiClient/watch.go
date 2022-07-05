package apiClient

import (
	"fmt"
	"log"

	"github.com/turbot/terraform-provider-turbot/errors"
)

func (client *Client) CreateWatch(input map[string]interface{}) (*Watch, error) {
	query := createWatchMutation()
	responseData := &WatchResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "watch")
	}
	log.Printf("Watch created: %s", responseData.Watch.Turbot.Id)

	return &responseData.Watch, nil
}

func (client *Client) WatchExists(id string) (bool, error) {
	resource, err := client.ReadWatch(id)

	if err != nil {
		if errors.NotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	exists := resource.Turbot.Id != ""
	return exists, nil
}

func (client *Client) ReadWatch(id string) (*Watch, error) {
	query := readWatchQuery(id)
	var responseData = &WatchResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "watch")
	}

	return &responseData.Watch, nil
}

func (client *Client) UpdateWatch(input map[string]interface{}) (*Watch, error) {
	query := updateWatchMutation()
	responseData := &WatchResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "watch")
	}
	return &responseData.Watch, nil
}

func (client *Client) DeleteWatch(id string) error {
	log.Printf("Deleting watch: %s", id)
	query := deleteWatchMutation(id)
	var responseData interface{}

	// execute api call
	if err := client.doRequest(query, nil, &responseData); err != nil {
		return fmt.Errorf("error deleting grant: %s", err.Error())
	}
	log.Printf("Watch deleted: %s", id)
	return nil
}
