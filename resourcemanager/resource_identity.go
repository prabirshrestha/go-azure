package resourcemanager

type ResourceIdentity struct {
	ResourceName               string
	ResourceType               string
	ResourceProviderApiVersion string
	ResourceProviderNamespace  string
	ParentResourcePath         string
}
