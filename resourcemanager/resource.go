package resourcemanager

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

}

type ResourceCreateOrUpdateResult struct {

}

type ResourceGetResult struct {
	
}

type Resource struct {
	Id       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

type ResourceListResult struct {
	Value []Resource `json:"value"`
}

func (ro *ResourceOperations) List(parameters *ResourceListParameters) (*ResourceListResult, *AzureOperationResponse, error) {
	subscriptionId := getSubscriptionId(ro.c, nil)

	var result ResourceListResult
	azureOperationResponse, err := ro.c.DoGet("/subscriptions/"+subscriptionId+"/resources?api-version="+ro.c.apiVersion, &result)

	if err != nil {
		return nil, nil, err
	}

	return &result, azureOperationResponse, nil
}

func (ro *ResourceOperations) ListNext(nextLink string) (*ResourceListResult, *AzureOperationResponse, error) {
	return nil, nil, nil
}

func (ro *ResourceOperations) Moveresources(sourceResourceGroupName string, parameters *ResourceMoveInfo) (*AzureOperationResponse, error) {
	return nil, nil
}

func (ro *ResourceOperations) Get(resourceGroupName string, identity *ResourceIdentity) (*ResourceGetResult, *AzureOperationResponse, error) {
	return nil, nil, nil
}

func (ro *ResourceOperations) Delete(resourceGroupName string, identity *ResourceIdentity) (*AzureOperationResponse, error) {
	return nil, nil
}

func (ro *ResourceOperations) CreateOrUpdate(resourceGroupName string, identity *ResourceIdentity) (*ResourceCreateOrUpdateResult, *AzureOperationResponse, error) {
	return nil, nil, nil
}

func (ro *ResourceOperations) CheckExistence(resourceGroupName string, identity *ResourceIdentity) (*ResourceExistsResult, *AzureOperationResponse, error) {
	return nil, nil, nil
}