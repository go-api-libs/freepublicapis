package main

import "net/http"

// TODO: generate from openapi.json
const defaultServerURL = "https://www.freepublicapis.com/api/"

func probe() error {
	_, err := http.Get(defaultServerURL + "random")
	return err
}
