// Package ctxval provides type-safe, collision-free accessors for values
// stored in a [context.Context] using Go generics.
//
// # Motivation
//
// The standard [context.WithValue] / [context.Context.Value] API is
// intentionally untyped: both the key and the value are of type any.
// This forces every caller to perform a type assertion and opens the door
// to key collisions between packages that happen to use the same key value
// (especially string keys).
//
// ctxval eliminates both problems by deriving the context key from the
// value type itself:
//
//	type key[T any] struct{}
//
// Each distinct type T yields a unique key[T]{}, so two packages can never
// collide unless they deliberately share the exact same type. No manual
// key declaration is required.
//
// # Basic usage
//
//	ctx = ctxval.With(ctx, "request-id-123")
//	id, ok := ctxval.Get[string](ctx)
//
//	user := ctxval.Must[*User](ctx)          // panics if absent
//	timeout := ctxval.GetOr(ctx, 30*time.Second)
//
// # Semantics
//
// The package follows the exact semantics of the standard library:
//
//   - With returns a derived context; the original context is never mutated.
//   - Subsequent calls with the same type T shadow earlier values; the full
//     history remains reachable through the parent chain.
//   - Get returns the zero value of T and false when the type has never been
//     stored.
//   - A typed nil (for example a nil *User) is stored and retrieved successfully
//     (ok == true). An untyped nil interface is not.
//
// # When to use Must
//
// Must is intended for values that are guaranteed to be present by upstream
// middleware (authentication, request-id injection, tenant resolution, …).
// Its absence is treated as an unrecoverable programming error. Prefer Get
// or GetOr for optional / best-effort data.
//
// # Limitation
//
// Because the key is derived solely from the type, only one value of a given
// type can be active in a context chain at a time (the most recent one).
// If you need multiple independent values of the same concrete type, define
// distinct named types or use an explicit-key library.
package ctxval
