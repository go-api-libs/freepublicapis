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
	c, err := freepublicapis.NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
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
		Limit: 350,
		Sort:  freepublicapis.SortBest,
	}); err != nil {
		return err
	}

	return nil
}
