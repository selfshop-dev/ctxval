# ctxval

Type-safe helpers for storing request-scoped values in a `context.Context`
without declaring a key for every value.

`context.WithValue` requires an unexported key type per value, which in practice
means a package-level key and a `With`/`From` pair for every piece of data.
`ctxval` uses the value's static type to instantiate its key type instead, so
`With[T]` and `Get[T]` work for any type without declaring keys at the call
site.

## Features

- one slot per Go type, with no exported key variables and no caller-side type
  assertions;
- distinct named types with the same underlying type stay independent
  (`userID int` and `teamID int` do not collide);
- stored zero values, nil pointers, and nil interface values are reported as
  present, so they remain distinguishable from missing data;
- derived contexts preserve cancellation, deadlines, and unrelated values.

## Installation

```sh
go get github.com/selfshop-dev/ctxval
```

The module requires Go 1.26.10 or newer, as declared by the `go` directive in
`go.mod`.

## Quick start

```go
package main

import (
	"context"
	"fmt"

	"github.com/selfshop-dev/ctxval"
)

type (
	userID    int
	requestID string
)

func handle(ctx context.Context) {
	if id, ok := ctxval.Get[userID](ctx); ok {
		fmt.Println("user:", id)
	}
}

func main() {
	ctx := ctxval.With(context.Background(), userID(42))
	ctx = ctxval.With(ctx, requestID("req-001"))

	handle(ctx)

	reqID, _ := ctxval.Get[requestID](ctx)
	fmt.Println("request:", reqID)
}
```

Output:

```text
user: 42
request: req-001
```

## Behavior

- `With` derives a new context; the parent context is never modified.
- Storing the same type again in a child context overrides that type only.
  Values of other types are inherited from the parent.
- `Get` returns the zero value of `T` and `false` when no value is present.
- A stored zero value or nil is reported with `ok == true`, so it is not
  confused with a missing value.
- Values are stored using ordinary Go assignment and are not deep-copied.
  Pointers, maps, and slice backing arrays may refer to data shared with the
  caller and other contexts. The package does not synchronize access to mutable
  data, so callers are responsible for safe concurrent access.
- Stored values remain retrievable after the context is cancelled.
- `With` panics if the parent context is nil, and `Get` panics if the context is
  nil.

## Limitations

- Only one value per type is visible in a context chain.
- The key is derived from the static type argument, not from the dynamic type
  of the value. Storing a value as `any` or `error` makes it retrievable only
  with that interface type.
- Two packages that store the same type in one context share a slot. Libraries
  should define dedicated named types instead of reusing `string`, `int`, or
  other common types.
- An untyped `nil` cannot be passed to `With`, because `T` cannot be inferred;
  use a typed nil or an explicit type argument such as `ctxval.With[error](ctx,
  nil)`.

A context is not a replacement for explicit function parameters. Prefer passing
data directly when the call chain allows it, and use `ctxval` for values that
must cross interfaces whose signatures cannot carry them, such as middleware and
request handlers.

## Documentation and examples

API documentation is available on
[pkg.go.dev](https://pkg.go.dev/github.com/selfshop-dev/ctxval).
The runnable package example lives in [example_test.go](example_test.go); it
has checked output and can be executed with:

```sh
go test -run Example ./...
```

## Development

```sh
./bin/golangci-lint fmt
./bin/golangci-lint run ./...
go test ./...
go test -race ./...
```

Run `go mod tidy` when dependencies or module metadata change. CI runs the
tests with the race detector and requires 100% statement coverage.

See [TOOLS.md](TOOLS.md) for installing the development tools.

## License

`ctxval` is released under the [MIT License](LICENSE).
