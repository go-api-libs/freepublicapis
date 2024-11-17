package main

import "net/http"

// TODO: generate from openapi.json
const defaultServerURL = "https://www.freepublicapis.com/api/"

// probe calls the API server to check what we can do
func probe() error {
	// define http calls here, e.g.: http.Get(defaultServerURL + "my-endpoint")

	_, err := http.Get(defaultServerURL + "random")
	return err
}
