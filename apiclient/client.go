package apiclient

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/machinebox/graphql"
	"net/url"
	"path"
)

// Turbot API Client
type Client struct {
	AccessKeyId, SecretAccessKey string
	Graphql                      *graphql.Client
}

func CreateClient(accessKeyId, secretAccessKey, workspace string) (*Client, error) {
	// build api url
	u, err := url.Parse(workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to create client - could not parse workspace url %s, error %s", workspace, err.Error())
	}
	if u.Path == "invalid" {
		return nil, fmt.Errorf("failed to create client - could not parse workspace url '%s'", workspace)
	}
	u.Path = path.Join(u.Path, "api/v5/graphql")
	baseUrl := u.String()

	return &Client{
		AccessKeyId:     accessKeyId,
		SecretAccessKey: secretAccessKey,
		Graphql:         graphql.NewClient(baseUrl),
	}, nil
}

// Validate checks if the API workspace URL and credentials are valid.
func (client *Client) Validate() error {
	query, responseObject := validationQuery()
	err := client.doRequest(query, nil, &responseObject)
	if err == nil && !responseObject.isValid() {
		err = errors.New("authorisation failed - have access_key_id and secret_access_key been set correctly?")
	}
	return err
}

func basicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// execute graphql request
func (client *Client) doRequest(query string, vars map[string]interface{}, responseData interface{}) error {
	// make a request
	req := graphql.NewRequest(query)

	// set any variables
	if vars != nil {
		for k, v := range vars {
			req.Var(k, v)
		}
	}

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", basicAuthHeader(client.AccessKeyId, client.SecretAccessKey))

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	if err := client.Graphql.Run(ctx, req, &responseData); err != nil {
		return err
	}
	return nil
}
