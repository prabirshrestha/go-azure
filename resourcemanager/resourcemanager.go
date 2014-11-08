package resourcemanager

import (
	"errors"
	"net/http"
)

const (
	defaultBasePath   = "https://management.azure.com/"
	defaultApiVersion = "2014-04-01-preview"
)

func New(options *Options) (*Service, error) {
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

	return s, nil
}

type Options struct {
	Client     *http.Client
	BasePath   string
	ApiVersion string
}

type Service struct {
	client     *http.Client
	BasePath   string
	ApiVersion string
}
