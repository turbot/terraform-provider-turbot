package apiClient

import "time"

type ClientConfig struct {
	Credentials     ClientCredentials
	CredentialsPath string
	Profile         string
	// RequestTimeout bounds every GraphQL request. Zero leaves it to CreateClient, which
	// installs DefaultRequestTimeout.
	RequestTimeout time.Duration
}

type ClientCredentials struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Workspace string
}
