package resourcemanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/prabirshrestha/go-azure/azure"
)

const (
	defaultBasePath   = "https://management.azure.com/"
	defaultApiVersion = "2015-01-01"
)

type Options struct {
	Client     *http.Client
	BasePath   string
	ApiVersion string

	Credentials interface{}
}

type ResourceManagementClient struct {
	client     *http.Client
	basePath   string
	apiVersion string

	tokenCredentials       *azure.TokenCredentials
	certificateCredentials *azure.CertificateCredentials

	Resource *ResourceOperations
}

func New(options *Options) (*ResourceManagementClient, error) {
	basePath := options.BasePath
	if basePath == "" {
		basePath = defaultBasePath
	}

	apiVersion := options.ApiVersion
	if apiVersion == "" {
		apiVersion = defaultApiVersion
	}

	client := &ResourceManagementClient{
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

// Generates a HTTP request but does not perform the request
func (c *ResourceManagementClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	var rbody io.Reader
	var contentType string

	switch t := body.(type) {
	case nil:
	case string:
		rbody = bytes.NewBufferString(t)
	case io.Reader:
		rbody = t
	default:
		v := reflect.ValueOf(body)
		if !v.IsValid() {
			break
		}

		if v.Type().Kind() == reflect.Ptr {
			v = reflect.Indirect(v)
			if !v.IsValid() {
				break
			}
		}

		j, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		rbody = bytes.NewReader(j)
		contentType = "application/json"
	}

	apiURL := strings.TrimRight(c.basePath, "/")
	if apiURL == "" {
		apiURL = defaultBasePath
	}

	req, err := http.NewRequest(method, apiURL+path, rbody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content=Type", contentType)
	}

	return req, nil
}

func (c *ResourceManagementClient) Do(request *http.Request, v interface{}) error {
	httpClient := c.client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return errors.New("error occurred")
	}

	switch t := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(t, res.Body)
	default:
		err = json.NewDecoder(res.Body).Decode(v)
	}

	return err
}

func getSubscriptionId(c *ResourceManagementClient, options interface{}) string {
	var result string

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

	return result
}
