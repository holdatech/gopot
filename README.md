[![GoDoc](https://godoc.org/github.com/holdatech/gopot/v2?status.svg)](https://pkg.go.dev/github.com/holdatech/gopot/v4)
[![Go](https://github.com/holdatech/gopot/workflows/Go/badge.svg)](https://github.com/holdatech/gopot/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/holdatech/gopot)](https://goreportcard.com/report/github.com/holdatech/gopot)

# gopot
Platform of Trust utility functions for go.

## Getting the library

```
go get -u github.com/holdatech/gopot/v4
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/holdatech/gopot/v4"
)

func main() {
	var secretKey, _ = rsa.GenerateKey(rand.Reader, 4096)

	payload := struct {
		Hello string `json:"hello"`
	}{
		Hello: "World",
	}

	signature, err := gopot.CreateSignature(payload, secretKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(signature)
}

```
