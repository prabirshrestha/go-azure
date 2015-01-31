package resourcemanager

import (
	"errors"
)

func NewResourceOperations(c *ResourceManagementClient) *ResourceOperations {
	ro := &ResourceOperations{c: c}
	return ro
}

type ResourceOperations struct {
	c *ResourceManagementClient
}

type ResourceListParameters struct {
	SubscriptionId string
	Top            int
	SkipToken      string
	Filter         string
	ApiVersion     string
}

type ResourceListResult struct {

}

func (ro *ResourceOperations) List(parameters ResourceListParameters) (*ResourceListResult, error) {
	subscriptionId := getSubscriptionId(ro.c, parameters)

	if subscriptionId == "" {
		return nil, errors.New("subscriptionId is empty")
	}

	return &ResourceListResult{}, nil
}
