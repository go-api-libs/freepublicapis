package main

import (
	"context"
	"fmt"

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

	if _, err := c.GetAPI(ctx, 1); err == nil {
		return fmt.Errorf("expected error")
	} else {
		fmt.Printf("err: %v\n", err)
	}

	if _, err := c.ListApis(ctx, nil); err != nil {
		return err
	}

	if _, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
		Limit: 10,
		Sort:  freepublicapis.SortBest,
	}); err != nil {
		return err
	}

	// TODO:
	// - enable 500 error but ignore content type (not important)
	// - ReadMe generation
	// - next APIs
	// - script that checks previous APIs, generates all 350 APIs

	return nil
}
