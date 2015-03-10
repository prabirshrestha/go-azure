package resourcemanager

import (
	"fmt"
	"net/url"
	"strconv"
)

func NewResourceOperations(c *ResourceManagementClient) *ResourceOperations {
	ro := &ResourceOperations{c: c}
	return ro
}

type ResourceOperations struct {
	c *ResourceManagementClient
}

type ResourceListParameters struct {
	ResourceGroupName string
	ResourceType      string
	TagName           string
	TagValue          string
	Top               int
}

type ResourceMoveInfo struct {
}

type ResourceExistsResult struct {
	Exists bool
}

type ResourceCreateOrUpdateResult struct {
}

type ResourceGetResult struct {
	Resource
}

type Resource struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Location   string            `json:"location"`
	Tags       map[string]string `json:"tags"`
	Properties interface{}       `json:"properties"`
}

type ResourceListResult struct {
	Value []Resource `json:"value"`
	Next  string     `json:"next"`
}

type ResourcesMoveInfo struct {
	Resources           []string `json:"resources"`
	TargetResourceGroup string   `json:"targetResourceGroup"`
}

// List all resources for a given Subscription and Resource group. Each resource is of type Resource.
func (ro *ResourceOperations) List(parameters *ResourceListParameters) (*ResourceListResult, *AzureOperationResponse, error) {
	subscriptionId := getSubscriptionId(ro.c, nil)

	path := "/subscriptions/" + url.QueryEscape(subscriptionId)

	if parameters != nil {
		if parameters.ResourceGroupName != "" {
			path += "/resourcegroups/" + url.QueryEscape(parameters.ResourceGroupName)
		}
	}

	path += "/resources?api-version=" + url.QueryEscape(ro.c.apiVersion)

	if parameters != nil {
		if parameters.Top != 0 {
			path += "&$top=" + strconv.Itoa(parameters.Top)
		}

		filter := ""

		if parameters.ResourceType != "" {
			filter += url.QueryEscape("resourceType eq '" + parameters.ResourceType + "'")
		}

		if parameters.TagValue != "" {
			if filter != "" {
				filter += url.QueryEscape(" and ")
			}
			filter += url.QueryEscape("tagValue eq '" + parameters.ResourceType + "'")
		}

		if filter != "" {
			path += "&filter=" + filter
		}
	}

	var result ResourceListResult
	azureOperationResponse, err := ro.c.DoGet(path, &result)

	if err != nil {
		return nil, nil, err
	}

	return &result, azureOperationResponse, nil
}

func (ro *ResourceOperations) ListNext(nextLink string) (*ResourceListResult, *AzureOperationResponse, error) {
	return nil, nil, nil
}

// Move resources from one resource group to another
func (ro *ResourceOperations) MoveResources(sourceResourceGroupName string, parameters *ResourcesMoveInfo) (*AzureOperationResponse, error) {
	subscriptionId := getSubscriptionId(ro.c, nil)

	path := fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/moveresources?api-version=%s",
		subscriptionId, sourceResourceGroupName, ro.c.apiVersion)

	return ro.c.DoPost(path, parameters, nil)
}

// Get a resource group
func (ro *ResourceOperations) Get(resourceGroupName string, identity *ResourceIdentity) (*ResourceGetResult, *AzureOperationResponse, error) {
	subscriptionId := getSubscriptionId(ro.c, nil)

	var result ResourceGetResult

	path := fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/%s/%s/%s/%s?api-version=%s",
		subscriptionId, resourceGroupName, identity.ResourceProviderNamespace, identity.ParentResourcePath, identity.ResourceType, identity.ResourceName,
		identity.ResourceProviderApiVersion)

	azureOperationResponse, err := ro.c.DoGet(path, &result)

	if err != nil {
		return nil, nil, err
	}

	return &result, azureOperationResponse, nil
}

// Delete a resource group
func (ro *ResourceOperations) Delete(resourceGroupName string, identity *ResourceIdentity) (*AzureOperationResponse, error) {
	subscriptionId := getSubscriptionId(ro.c, nil)

	path := fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/%s/%s/%s/%s?api-version=%s",
		subscriptionId, resourceGroupName, identity.ResourceProviderNamespace, identity.ParentResourcePath, identity.ResourceType, identity.ResourceName,
		identity.ResourceProviderApiVersion)

	return ro.c.DoDelete(path, nil)
}

// Create a resource group.
func (ro *ResourceOperations) CreateOrUpdate(resourceGroupName string, identity *ResourceIdentity) (*ResourceCreateOrUpdateResult, *AzureOperationResponse, error) {
	return nil, nil, nil
}

// Checks whether resource group exists.
func (ro *ResourceOperations) CheckExistence(resourceGroupName string, identity *ResourceIdentity) (*ResourceExistsResult, *AzureOperationResponse, error) {
	_, azureOperationResponse, err := ro.Get(resourceGroupName, identity)

	result := ResourceExistsResult{Exists: true}

	if err != nil {
		switch err.(type) {
		case Error:
			error := err.(Error)
			if error.StatusCode == 404 {
				result.Exists = false
				return &result, error.AzureOperationResponse, nil
			} else {
				return nil, nil, err
			}
		default:
			return nil, nil, err
		}
	}

	return &result, azureOperationResponse, nil
}
