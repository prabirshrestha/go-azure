package main

import "github.com/prabirshrestha/go-azure/azure"
import arm "github.com/prabirshrestha/go-azure/resourcemanager"
import "fmt"

func main() {
	token, _ := azure.NewTokenCredentials("subscriptionId-a", "token-a")
	client, _ := arm.New(arm.Options{Credentials: token})

	parameters := arm.ResourceListParameters{}
	result, _ := client.Resource.List(parameters)

	fmt.Println(client)
	fmt.Println(result)
}