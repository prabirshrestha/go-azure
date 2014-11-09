package resourcemanager

import (
	"errors"
	"net/http"
)

const (
	defaultBasePath   = "https://management.azure.com/"
	defaultApiVersion = "2014-04-01-preview"
)

func NewWithToken(token string) (*Service, error) {
	if token == "" {
		return nil, errors.New("token in empty")
	}

	options := &Opts{
		Client: &http.Client{},
	}

	return New(options)
}

func New(options *Opts) (*Service, error) {
	if options == nil {
		return nil, errors.New("options is nil")
	}

	if options.Client == nil {
		return nil, errors.New("options.Client is nil")
	}

	basePath := options.BasePath
	if basePath == "" {
		basePath = defaultBasePath
	}

	apiVersion := options.ApiVersion
	if apiVersion == "" {
		apiVersion = defaultApiVersion
	}

	s := &Service{
		client:     options.Client,
		BasePath:   basePath,
		ApiVersion: apiVersion,
	}

	s.Resource = NewResourceService(s)

	return s, nil
}

type Opts struct {
	Client     *http.Client
	BasePath   string
	ApiVersion string
}

type Service struct {
	client     *http.Client
	BasePath   string
	ApiVersion string

	Resource *ResourceService
}
