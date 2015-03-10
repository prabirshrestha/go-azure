package resourcemanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/prabirshrestha/go-azure/azure"
)

const (
	defaultBasePath   = "https://management.azure.com/"
	defaultApiVersion = "2015-01-01"
)

// Options to pass in while creating ResourceManager Client
type Options struct {
	Client     *http.Client
	BasePath   string
	ApiVersion string

	Credentials interface{}
}

// Client to access ResourceManager APIs
type ResourceManagementClient struct {
	client     *http.Client
	basePath   string
	apiVersion string

	tokenCredentials       *azure.TokenCredentials
	certificateCredentials *azure.CertificateCredentials

	Resources *ResourceOperations
}

// Base Response object
type AzureOperationResponse struct {
	RequestId                           string
	RateLimitRemainingSubscriptionReads int
	CorrelationRequestId                string
	RoutingRequestId                    string
	StatusCode                          int
}

type Error struct {
	error
	AzureOperationResponse *AzureOperationResponse
	StatusCode             int
	Body                   interface{}
}

// Create a new client
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

	client.Resources = NewResourceOperations(client)

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
		req.Header.Set("ContentType", contentType)
	}

	if c.tokenCredentials != nil {
		req.Header.Set("Authorization", "Bearer "+c.tokenCredentials.Token)
	}

	return req, nil
}

// Perform an operation on the API
func (c *ResourceManagementClient) Do(request *http.Request, result interface{}) (*AzureOperationResponse, error) {
	httpClient := c.client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	azureOperationResponse := &AzureOperationResponse{}
	azureOperationResponse.RateLimitRemainingSubscriptionReads, _ = strconv.Atoi(res.Header.Get("x-ms-ratelimit-remaining-subscription-reads"))
	azureOperationResponse.RequestId = res.Header.Get("x-ms-request-id")
	azureOperationResponse.CorrelationRequestId = res.Header.Get("x-ms-correlation-request-id")
	azureOperationResponse.RoutingRequestId = res.Header.Get("x-ms-routing-request-id")
	azureOperationResponse.StatusCode = res.StatusCode

	if res.StatusCode/100 != 2 {

		error := Error{
			AzureOperationResponse: azureOperationResponse,
			StatusCode:             res.StatusCode,
		}

		contentType := res.Header.Get("content-type")

		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			error.error = err
			return nil, error
		} else {
			if strings.Contains(contentType, "application/json") {
				err := json.Unmarshal(contents, &error.Body)
				if err != nil {
					error.error = err
					return nil, error
				} else {
					error.error = errors.New(fmt.Sprintf("go-azure error: Status Code: %v Body: %v", res.StatusCode, error.Body))
					return nil, error
				}
			} else {
				error.Body = string(contents)
				error.error = errors.New(fmt.Sprintf("go-azure error: Status Code: %v Body: %v", res.StatusCode, error.Body))
				return nil, error
			}
		}
	}

	switch t := result.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(t, res.Body)
	default:
		err = json.NewDecoder(res.Body).Decode(result)
	}

	if err != nil {
		return nil, err
	}

	return azureOperationResponse, nil
}

// Perform a GET operation
func (c *ResourceManagementClient) DoGet(path string, result interface{}) (*AzureOperationResponse, error) {
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, result)
}

// Perform a POST operation
func (c *ResourceManagementClient) DoPost(path string, body interface{}, result interface{}) (*AzureOperationResponse, error) {
	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req, result)
}

// Perform a PUT operation
func (c *ResourceManagementClient) DoPut(path string, body interface{}, result interface{}) (*AzureOperationResponse, error) {
	req, err := c.NewRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req, result)
}

// Perform a PATCH operation
func (c *ResourceManagementClient) DoPatch(path string, body interface{}, result interface{}) (*AzureOperationResponse, error) {
	req, err := c.NewRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req, result)
}

// Perform a DELETE operation
func (c *ResourceManagementClient) DoDelete(path string, result interface{}) (*AzureOperationResponse, error) {
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, result)
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
