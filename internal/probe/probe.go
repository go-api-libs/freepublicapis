package main

import (
	"context"

	"github.com/go-api-libs/freepublicapis/pkg/freepublicapis"
)

// TODO: generate from openapi.json
// const defaultServerURL = "https://www.freepublicapis.com/api/"

// probe calls the API server to check what we can do
func probe() error {
	// define http calls here, e.g.: http.Get(defaultServerURL + "my-endpoint")

	ctx := context.Background()
	c, err := freepublicapis.NewClient()
	if err != nil {
		return err
	}

	if _, err := c.GetRandom(ctx); err != nil {
		return err
	}

	if _, err := c.GetAPI(ctx, 275); err != nil {
		return err
	}

	// _, err := http.Get("https://www.freepublicapis.com/api/apis/275")
	return nil
}
