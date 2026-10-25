// Package ctxval stores and retrieves request-scoped values in a
// [context.Context] using the value's static type as the key.
//
// # Overview
//
// The [context] package recommends unexported key types so that different
// packages cannot collide, which usually means declaring a key type and a
// With/From pair for every value. ctxval instead derives the key from the type
// argument: [With] and [Get] are generic over T and use key[T]{} as the context
// key. The key type key[T] is instantiated with each type argument, so callers
// do not declare context keys or perform type assertions. The package
// encapsulates both operations and keeps no package-level mutable state.
//
// # Core concepts
//
// [With] returns a derived context that carries a value under the key for T,
// and [Get] performs the matching lookup. Each distinct type has its own slot.
// Defined types with the same underlying type remain independent, so userID and
// teamID can both be stored without interfering. Type aliases share a slot with
// their target, because they denote the same type.
//
// At most one value per type is visible through a context. Storing a value of a
// type that already exists in an ancestor context overrides it for the derived
// context only; the parent context and its other values are unchanged.
//
// # Basic usage
//
//	type userID int
//
//	ctx := ctxval.With(parent, userID(42))
//	id, ok := ctxval.Get[userID](ctx)
//
// See the package example in example_test.go for a complete program that
// stores several types, derives a child context, and overrides a single value.
//
// # Behavior and guarantees
//
// [Get] reports a missing value as the zero value of T and false. Stored values
// are wrapped, so a stored zero value, nil pointer, or nil interface value is
// reported as present. Values remain retrievable after the context is cancelled,
// and because [With] derives the context with [context.WithValue], deadlines,
// cancellation, and unrelated values continue to work as usual.
//
// Stored values are assigned with ordinary Go assignment and are not
// deep-copied. Pointers, maps, and slice backing arrays may refer to data shared
// with the caller and with every context derived from the result. The package
// does not synchronize access to mutable data, so callers are responsible for
// safe concurrent access. Apart from that, the returned context is safe for
// simultaneous use by multiple goroutines, following the rules of
// [context.Context].
//
// # Limitations
//
// The key is derived from the static type argument, not from the dynamic type
// of the value.
// Storing a value as an interface type, such as error or any, makes it
// retrievable only with that interface type; [Get] with a concrete type does not
// find it. Conversely, two packages that store the same type in one context
// share a slot, so libraries should model their values as dedicated named types
// rather than reusing common types such as string or int.
//
// An untyped nil cannot be passed to [With], because T cannot be inferred from
// it; store a typed nil or pass an explicit type argument such as
// With[error](ctx, nil) when a nil value must be recorded.
//
// A context is not a substitute for explicit parameters: prefer passing data
// directly when the call chain allows it, and use ctxval for data that must
// cross interfaces whose signatures cannot carry it, such as middleware and
// request handlers.
//
// # See also
//
// [With] and [Get] document the individual contracts. The package example in
// example_test.go demonstrates storing several types, deriving a child context,
// and reading a missing value; it has checked output.
package ctxval
