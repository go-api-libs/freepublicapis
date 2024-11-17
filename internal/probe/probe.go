package main

import "net/http"

func probe() error {
	http.Get()

	return nil
}
