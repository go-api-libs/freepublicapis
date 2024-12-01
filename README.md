# 🤯 Free Public APIs API
[![Go Reference](https://pkg.go.dev/badge/github.com/go-api-libs/freepublicapis.svg)](https://pkg.go.dev/github.com/go-api-libs/freepublicapis/pkg/freepublicapis)
[![Official Documentation](https://img.shields.io/badge/docs-API-blue)](https://www.freepublicapis.com/api)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1-blue)](/api/openapi.json)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-api-libs/freepublicapis)](https://goreportcard.com/report/github.com/go-api-libs/freepublicapis)
![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)
![API Health](https://img.shields.io/badge/API_health-90%25-green)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

FreePublicAPIs now has a free public API! The Free Public APIs API provides up-to-date information on every free public API listed on freepublicapis.com. Perfect for developers, students and hobbyists. It is limited to 1000 requests a day, if you need more, feel free to reach out. Also if you build something with this API, make sure to add it as a showcase project here. ([Source](https://www.freepublicapis.com/api))

## Installation

To install the library, use the following command:

```shell
go get github.com/go-api-libs/freepublicapis/pkg/freepublicapis
```

## Usage

### Example 1: Get a Random API

```go
package main

import (
	"context"

	"github.com/go-api-libs/freepublicapis/pkg/freepublicapis"
)

func main() {
	c, err := freepublicapis.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	simpleAPIInfo, err := c.GetRandom(ctx)
	if err != nil {
		panic(err)
	}

	// Use simpleAPIInfo object
}

```

### Example 2: Get a Specific API by ID

```go
package main

import (
	"context"

	"github.com/go-api-libs/freepublicapis/pkg/freepublicapis"
)

func main() {
	c, err := freepublicapis.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	simpleAPIInfo, err := c.GetAPI(ctx, 275)
	if err != nil {
		panic(err)
	}

	// Use simpleAPIInfo object
}

```

### Example 3: List APIs with Parameters

```go
package main

import (
	"context"

	"github.com/go-api-libs/freepublicapis/pkg/freepublicapis"
)

func main() {
	c, err := freepublicapis.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	apiInfos, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
		Limit: 10,
		Sort:  "best",
	})
	if err != nil {
		panic(err)
	}

	// Use apiInfos array
}

```

## Additional Information

- [**Go Reference**](https://pkg.go.dev/github.com/go-api-libs/freepublicapis/pkg/freepublicapis): The Go reference documentation for the client package.
- [**Official Documentation**](https://www.freepublicapis.com/api): The official API documentation.
- [**OpenAPI Specification**](./api/openapi.json): The OpenAPI 3.1.0 specification.
- [**Go Report Card**](https://goreportcard.com/report/github.com/go-api-libs/freepublicapis): Check the code quality report.

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/go-api-libs/freepublicapis).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
