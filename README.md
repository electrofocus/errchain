# errchain

[![Go Reference](https://pkg.go.dev/badge/github.com/electrofocus/errchain.svg)](https://pkg.go.dev/github.com/electrofocus/errchain)

## About

Here's [Go](https://go.dev) package for errors chaining for further examining using the standard `errors.Is`. You can learn more about working with errors in Go in [this](https://go.dev/blog/go1.13-errors) article. Explore the [example](#example) below for more understanding.

This package uses [module version numbering](https://go.dev/doc/modules/version-numbers).


## Install
With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain run:

```
go get github.com/electrofocus/errchain
```

## Examples

Let's build new error from multiple errors and examine it with `errors.Is`:

```go
package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/electrofocus/errchain"
)

func main() {
	var (
		myErr = errors.New("my err")
		err   = errchain.New(myErr, io.EOF, os.ErrClosed, http.ErrHijacked)
	)

	if errors.Is(err, io.EOF) {
		fmt.Printf("here we have %q error\n", io.EOF)
	}

	if errors.Is(err, myErr) {
		fmt.Printf("and %q error\n", myErr)
	}

	if errors.Is(err, os.ErrClosed) {
		fmt.Printf("and %q error\n", os.ErrClosed)
	}

	if errors.Is(err, http.ErrHijacked) {
		fmt.Printf("and %q error,\n", http.ErrHijacked)
	}

	if !errors.Is(err, http.ErrAbortHandler) {
		fmt.Printf("but don't have %q error\n", http.ErrAbortHandler)
	}
}
```

Open above example in [The Go Playground](https://go.dev/play/p/yfPyoY_yVPi).
