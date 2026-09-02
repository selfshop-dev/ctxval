package ctxval_test

import (
	"context"
	"fmt"
	"time"

	"github.com/selfshop-dev/ctxval"
)

func Example() {
	ctx := context.Background()

	// Store values of different types.
	ctx = ctxval.With(ctx, "req-abc-123")
	ctx = ctxval.With(ctx, 42)
	ctx = ctxval.With(ctx, 30*time.Second)

	// Retrieve them in a type-safe way.
	id, ok := ctxval.Get[string](ctx)
	fmt.Println(id, ok)

	n, ok := ctxval.Get[int](ctx)
	fmt.Println(n, ok)

	// GetOr returns the stored value when present.
	timeout := ctxval.GetOr(ctx, 10*time.Second)
	fmt.Println(timeout)

	// GetOr returns the fallback when the type is absent.
	missing := ctxval.GetOr(ctx, true) // bool was never stored
	fmt.Println(missing)

	// Exists is a convenient predicate.
	fmt.Println(ctxval.Exists[string](ctx))
	fmt.Println(ctxval.Exists[bool](ctx))

	// Output:
	// req-abc-123 true
	// 42 true
	// 30s
	// true
	// true
	// false
}

func ExampleMust() {
	type userID string

	ctx := context.Background()
	ctx = ctxval.With(ctx, userID("u-42"))

	// Must is appropriate when middleware guarantees the value.
	id := ctxval.Must[userID](ctx)
	fmt.Println(id)

	// Output:
	// u-42
}
