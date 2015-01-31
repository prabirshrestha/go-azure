package resourcemanager

import (
	"net/http"

	"github.com/prabirshrestha/go-azure/azure"
)

const (
	defaultBasePath = "https://management.azure.com/"
	defaultApiVersion = "2015-01-01"
)

type Options struct {
	Client *http.Client
	BasePath string
	ApiVersion string

	Credentials interface{}
}

type ResourceManagementClient struct {
	client *http.Client;
	basePath string;
	apiVersion string;

	tokenCredentials *azure.TokenCredentials;
	certificateCredentials *azure.CertificateCredentials;

	Resource *ResourceOperations
}

func New(options Options) (*ResourceManagementClient, error) {
	httpClient := options.Client;
	if options.Client == nil {
		httpClient = http.DefaultClient
	}

	basePath := options.BasePath
	if basePath == "" {
		basePath = defaultBasePath
	}

	apiVersion := options.ApiVersion
	if apiVersion == "" {
		apiVersion = defaultApiVersion
	}

	client := &ResourceManagementClient{
		client:     httpClient,
		basePath:   basePath,
		apiVersion: apiVersion,
	}

	if tokenCredentials, ok := options.Credentials.(azure.TokenCredentials); ok {
		client.tokenCredentials = &tokenCredentials
	}

	if certificateCredentials, ok := options.Credentials.(azure.CertificateCredentials); ok {
		client.certificateCredentials = &certificateCredentials
	}

	client.Resource = NewResourceOperations(client)

	return client, nil
}

func getSubscriptionId(c *ResourceManagementClient, options interface{}) string {
	var result string;

	if options != nil {
		// use subscription Id from options
	}

	if result == "" {
		if c.tokenCredentials != nil {
			result = c.tokenCredentials.SubscriptionId
		}

		if result != "" && c.certificateCredentials != nil {
			result = c.certificateCredentials.SubscriptionId
		}
	}

	return result;
}