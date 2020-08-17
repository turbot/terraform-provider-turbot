package apiClient

import (
	"fmt"
)

var LdapDirectoryProperties = []interface{}{
	map[string]interface{}{"parent": "turbot.parentId"},
	"title",
	"description",
	"profileIdTemplate",
	"groupProfileIdTemplate",
	"url",
	"distinguishedName",
	"password",
	"base",
	"userObjectFilter",
	"disabledUserFilter",
	"userMatchFilter",
	"userSearchFilter",
	"userSearchAttributes",
	"groupObjectFilter",
	"groupSearchFilter",
	"groupSyncFilter",
	"userCanonicalNameAttribute",
	"userEmailAttribute",
	"userDisplayNameAttribute",
	"userGivenNameAttribute",
	"userFamilyNameAttribute",
	"tlsEnabled",
	"tlsServerCertificate",
	"groupMemberOfAttribute",
	"groupMembershipAttribute",
	"connectivityTestFilter",
	"rejectUnauthorized",
	"disabledGroupFilter",
}

func (client *Client) ReadLdapDirectory(id string) (*LdapDirectory, error) {
	// create a map of the properties we want the graphql query to return
	query := readResourceQuery(id, LdapDirectoryProperties)
	responseData := &LdapDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) CreateLdapDirectory(input map[string]interface{}) (*LdapDirectory, error) {
	query := createLdapDirectoryMutation(LdapDirectoryProperties)
	responseData := &LdapDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error creating local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateLdapDirectory(input map[string]interface{}) (*LdapDirectory, error) {
	query := updateLdapDirectoryMutation(LdapDirectoryProperties)
	responseData := &LdapDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error updating local directory: %s", err.Error())
	}
	return &responseData.Resource, nil
}

func (client *Client) DeleteLdapDirectory(aka string) error {
	query := deleteLdapDirectory()
	// we do not care about the response
	var responseData interface{}

	variables := map[string]interface{}{
		"input": map[string]string{
			"id": aka,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting resource: %s", err.Error())
	}
	return nil
}
