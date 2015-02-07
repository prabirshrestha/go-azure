package main

import (
	"fmt"
	"os"

	"github.com/prabirshrestha/go-azure/azure"
	arm "github.com/prabirshrestha/go-azure/resourcemanager"
)

func main() {
	token, _ := azure.NewTokenCredentials(os.Getenv("subscription"), os.Getenv("token"))
	client, _ := arm.New(&arm.Options{Credentials: token})

	// listResources(client)
	getResource(client)
}

func listResources(client *arm.ResourceManagementClient) {
	parameters := &arm.ResourceListParameters{}
	result, aor, err := client.Resources.List(parameters)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", aor)
	fmt.Printf("%+v\n", result)
}

func getResource(client *arm.ResourceManagementClient) {
	ri := &arm.ResourceIdentity{}
	ri.ResourceProviderNamespace = "Microsoft.Web"
	ri.ResourceType = "sites"
	ri.ResourceName = "pstestwebsite"
	ri.ResourceProviderApiVersion = "2014-04-01"

	resource, aor, err := client.Resources.Get("Default-Web-WestUS", ri)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", aor)
	fmt.Printf("%+v\n", resource)
	properties := resource.Properties.(map[string]interface{})
	fmt.Println(properties["adminEnabled"])
}
