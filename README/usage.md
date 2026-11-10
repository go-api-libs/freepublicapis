To install the library, use the following command:

```shell
go get github.com/go-api-libs/freepublicapis/pkg/freepublicapis
```


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

	// Use apiInfos slice
}

```
