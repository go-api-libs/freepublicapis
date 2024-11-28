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

	if _, err := c.ListApis(ctx, nil); err != nil {
		return err
	}

	limit := 0
	sort := ""
	apis, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
		Limit: &limit,
		Sort:  &sort,
	})
	if err != nil {
		return err
	}

	fmt.Printf("len(*apis): %v\n", len(*apis))

	// TODO: you can now add a limit parameter to the 'list all free public APIs' endpoint. You can also sort by: best, new, fast, popular, noerror, reliable, all.
	// _, err = http.Get("https://www.freepublicapis.com/api/apis?limit=10&sort=best")
	// return err
	return nil
}
