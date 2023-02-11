# errchain

❗️ This package is no longer needed if you are using Go version >= 1.20, since [in 1.20 standard errors package was extended](https://go.dev/doc/go1.20#errors) by adding [`Join`](https://pkg.go.dev/errors#Join) method, which provides functionality similar to that provided by errchain package.

[![Go Reference](https://pkg.go.dev/badge/github.com/electrofocus/errchain.svg)](https://pkg.go.dev/github.com/electrofocus/errchain)

## About

Here's [Go](https://go.dev) package for errors chaining for further examining using the standard `errors.Is`. You can learn more about working with errors in Go in [this](https://go.dev/blog/go1.13-errors) article. Explore [example](#examples) below for more understanding.

This package uses [module version numbering](https://go.dev/doc/modules/version-numbers).


## Install
With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain run:

```
go get github.com/electrofocus/errchain
```

## Examples

### Chain and examine
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

### Check error for compliance with one of expected errors
Moreover, non-obvious potential of `errchain` package is the ability to examine an error for compliance with one of expected ones.

Let's declare a `toy` function, as a result of which we expect an `error` (in fact, it will always return `io.EOF`):
```go
import "encoding/json"

func toy() error {
	return io.EOF
}
```

Sometimes you expect several different errors. In this case, to recognize the returned error, you need to do something like this:
```go
import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := toy(); err != nil &&
		(errors.Is(err, os.ErrClosed) ||
			errors.Is(err, io.EOF) ||
			errors.Is(err, http.ErrHijacked)) {

		fmt.Printf("got one of expected errors: %q\n", err)
	}
}
```

An equivalent check can be performed using `errors.Is` and `errchain` package's `New` functions, but in more concise and convinient way:
```go
import (
	"errors"
	"fmt"
	"net/http"
	"os"
	
	"github.com/electrofocus/errchain"
)

func main() {
	if err := toy(); errors.Is(errchain.New(
		os.ErrClosed,
		io.EOF,
		http.ErrHijacked,
	), err) {
		fmt.Printf("got one of expected errors: %q\n", err)
	}
}
```

Play with above example in [The Go Playground](https://go.dev/play/p/LwaRS9gB1Bn).
