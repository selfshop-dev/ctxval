package ctxval

import (
	"context"
	"fmt"
	"reflect"
)

// key is the unexported generic context key type.
// Each distinct type T produces a unique key[T]{}, eliminating collisions
// between packages without any manual key declaration.
type key[T any] struct{}

// With returns a derived context that carries value under the type T.
//
// The returned context is a shallow wrapper around the parent; the original
// ctx is never mutated. Subsequent calls that store another value of the
// same type T shadow the previous one, exactly as [context.WithValue] does.
// The earlier values remain reachable by walking the parent chain, so code
// that holds a reference to an intermediate context can still observe them.
//
// Because the key is derived solely from the type parameter, two packages
// can never collide unless they deliberately use the identical concrete
// type for unrelated purposes. No exported key constants or manual
// declarations are required.
//
// A typed nil (for example a nil pointer or nil interface value of a
// concrete type) is stored successfully and can later be retrieved with
// ok == true. Passing an untyped nil interface is not supported and will
// not behave as a stored value of type T.
func With[T any](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, key[T]{}, value)
}

// Get retrieves the value of type T previously stored in the context chain
// by a call to [With].
//
// If a value of type T is present (including a typed nil), Get returns that
// value and true. If no value of type T has ever been stored in this context
// or any of its parents, Get returns the zero value of T and false.
//
// The lookup follows the standard context parent chain, so a value stored
// on a parent remains visible to all derived children unless it has been
// shadowed by a more recent [With] call for the same type.
func Get[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)
	return v, ok
}

// GetOr is a convenience wrapper around [Get] that returns a caller-supplied
// fallback when the requested type is absent.
//
// If a value of type T exists in the context chain, that value is returned
// unchanged. Otherwise the provided fallback is returned. This is useful
// for optional configuration, default timeouts, or any situation where a
// sensible default is preferable to an explicit presence check.
//
// GetOr never panics and never allocates beyond what the underlying
// context implementation already does.
func GetOr[T any](ctx context.Context, fallback T) T {
	if v, ok := Get[T](ctx); ok {
		return v
	}
	return fallback
}

// Exists reports whether a value of type T is present anywhere in the
// context chain.
//
// It is equivalent to the boolean result of [Get] but avoids materialising
// the value when the caller only needs to know about presence. Exists
// returns true even when the stored value is a typed nil.
func Exists[T any](ctx context.Context) bool {
	_, ok := Get[T](ctx)
	return ok
}

// Must retrieves the value of type T from the context chain and panics if
// the value is not present.
//
// Must is intended exclusively for values that are contractually guaranteed
// to have been injected by upstream middleware or setup code (request IDs,
// authenticated users, tenants, loggers scoped to the request, and similar).
// Its absence indicates a programming error that should surface immediately
// rather than be handled as a routine runtime condition.
//
// Prefer [Get] or [GetOr] for any data that may legitimately be missing.
// The panic message includes the concrete type name to aid debugging.
func Must[T any](ctx context.Context) T {
	if v, ok := Get[T](ctx); ok {
		return v
	}
	panic(fmt.Sprintf(
		"ctxval.Must: value of type %s not found in context",
		reflect.TypeFor[T](),
	))
}
