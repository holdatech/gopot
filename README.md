[![GoDoc](https://godoc.org/github.com/holdatech/gopot/v2?status.svg)](https://pkg.go.dev/github.com/holdatech/gopot/v2)
[![Go](https://github.com/holdatech/gopot/workflows/Go/badge.svg)](https://github.com/holdatech/gopot/actions)

# gopot
Platform of Trust utility functions for go.

## Usage

```go
package main

import (
	"fmt"

	"github.com/holdatech/gopot/v2"
)

func main() {
	secret := []byte("secret123")

	payload := struct {
		Hello string `json:"hello"`
	}{
		Hello: "World",
	}

	signature, err := gopot.CreateSignature(payload, secret)
	if err != nil {
		panic(err)
	}

	fmt.Println(signature)
}

```
