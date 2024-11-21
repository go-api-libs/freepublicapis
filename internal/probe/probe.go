package main

import (
	"context"
	"fmt"

	"github.com/go-api-libs/freepublicapis/pkg/freepublicapis"
)

// TODO: generate from openapi.json
const defaultServerURL = "https://www.freepublicapis.com/api/"

// probe calls the API server to check what we can do
func probe() error {
	// define http calls here, e.g.: http.Get(defaultServerURL + "my-endpoint")

	c, err := freepublicapis.NewClient()
	if err != nil {
		return err
	}

	api, err := c.GetRandom(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("api: %v\n", api)

	return nil

	// _, err := http.Get(defaultServerURL + "random")
	// return err
}
