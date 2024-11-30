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

```bash
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
	randomAPI, err := c.GetRandom(ctx)
	if err != nil {
		panic(err)
	}

	// Use randomAPI
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
	myAPI, err := c.GetAPI(ctx, 275)
	if err != nil {
		panic(err)
	}

	// Use myAPI
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
	apis, err := c.ListApis(ctx, &freepublicapis.ListApisParams{
		Limit: 100,
		Sort:  freepublicapis.SortBest,
	})
	if err != nil {
		panic(err)
	}

	// Use the apis
}
```

## Additional Information

- **Go Reference**: The Go reference documentation for this package can be found [here](https://pkg.go.dev/github.com/go-api-libs/freepublicapis/pkg/freepublicapis).
- **Official Documentation**: The official API documentation is available [here](https://www.freepublicapis.com/api).
- **OpenAPI Specification**: You can find the OpenAPI specification [here](/api/openapi.json).
- **Go Report Card**: Check the code quality report [here](https://goreportcard.com/report/github.com/go-api-libs/freepublicapis).

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/go-api-libs/freepublicapis).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.




