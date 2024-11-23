package main

import (
	"net/http"
)

// TODO: generate from openapi.json
// const defaultServerURL = "https://www.freepublicapis.com/api/"

// probe calls the API server to check what we can do
func probe() error {
	// define http calls here, e.g.: http.Get(defaultServerURL + "my-endpoint")

	// c, err := freepublicapis.NewClient()
	// if err != nil {
	// 	return err
	// }

	// api, err := c.GetRandom(context.Background())
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("API Description: %v\n", api.Description)

	_, err := http.Get("https://www.freepublicapis.com/api/apis/275")
	return err
}
