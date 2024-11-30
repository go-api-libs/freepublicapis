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

	c, err := freepublicapis.NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	randomAPI, err := c.GetRandom(ctx)
	if err != nil {
		return err
	}
	_ = randomAPI

	myAPI, err := c.GetAPI(ctx, 275)
	if err != nil {
		return err
	}

	fmt.Printf("api: %#v\n", myAPI)

	if _, err := c.GetAPI(ctx, 1); err == nil {
		return fmt.Errorf("expected error")
	} else {
		fmt.Printf("err: %v\n", err)
	}

	if _, err := c.ListApis(ctx, nil); err != nil {
		return err
	}

	apis, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
		Limit: 1000,
		Sort:  freepublicapis.SortBest,
	})
	if err != nil {
		return err
	}

	for _, api := range apis {
		if api.Documentation.String() != api.Source.String() {
			fmt.Printf("%v\nDoc: %s\nSource: %s\n", api.Title, &api.Documentation, &api.Source)
		}
	}

	// TODO:
	// - ReadMe generation (use ReadMes from other projects as example)
	// - next APIs
	// - script that checks previous APIs, generates all 350 APIs

	return nil
}
