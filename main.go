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

	parameters := &arm.ResourceListParameters{}
	result, opResponse, err := client.Resource.List(parameters)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", opResponse)
	fmt.Printf("%+v\n", result)
}
