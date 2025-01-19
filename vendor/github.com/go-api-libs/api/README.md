# Standardized API Error Handling for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/go-api-libs/api.svg)](https://pkg.go.dev/github.com/go-api-libs/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-api-libs/api)](https://goreportcard.com/report/github.com/go-api-libs/api)
![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

This library provides a set of common error definitions used across various API interactions. By using this library, you can achieve consistent error handling in your applications, making it easier to debug, maintain, and integrate with different services.

It is utilized by all API libraries within the [go-api-libs](https://github.com/go-api-libs) organization.

## Features

- **Wrap Errors**: Easily wrap underlying errors with API-specific errors to provide more context.
- **Error Identification**: Check and identify the nature of errors using Go's `errors.Is` and `errors.As` functions.
- **Custom Errors**: Define and use custom errors specific to your API needs.
- **Status Code Handling**: Handle and identify errors based on HTTP status codes.
- **Content Type Handling**: Handle errors based on content type mismatches.

## Installation

To install the library, run:
```sh
go get github.com/go-api-libs/api
```

## Usage

Here is a basic example of how to use the api library for error handling:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-api-libs/api"
	"github.com/go-api-libs/toggl/pkg/toggl"
)

func main() {
	c, err := toggl.NewClientWithAPIToken(os.Getenv("TOGGL_TOKEN"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	now := time.Now()
	entries, err := c.ListTimeEntriesInRange(ctx, now.Add(-24*time.Hour), now)
	if err != nil {
		decErr := &api.DecodingError{}
		if errors.As(err, &decErr) {
			fmt.Println("Could not decode response:", decErr.Err)
		}

		apiErr := &api.Error{}
		if errors.As(err, &apiErr) {
			fmt.Println("Response Status Code:", apiErr.StatusCode())
			fmt.Println("Response Content Type:", apiErr.ContentType())
			if apiErr.IsCustom {
				fmt.Println("API sent back the error:", apiErr.Err)
				return
			}
		}

		if errors.Is(err, api.ErrStatusCode) {
			fmt.Println("The response status code indicates an error (but is properly documented).")
			// NOTE: Some APIs have custom error responses that don't contain api.ErrStatusCode.
			// If the status code indicates an error and the API returns a JSON,
			// we unmarshal it and return it as an object that fulfils the error interface.
			// If this happens, `IsCustom` above is set to true.
		} else if errors.Is(err, api.ErrUnknownStatusCode) {
			fmt.Println("The returned status code is not documented in the API specification.")
		}

		if errors.Is(err, api.ErrUnknownContentType) {
			fmt.Println("The response content type is not documented in the API specification.")
		}

		panic(err)
	}

	// Use entries slice
}
```

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/go-api-libs/api).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
