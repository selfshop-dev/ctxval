package ctxval

import "context"

type (
	// key is the context key for values of type T. Because the key is an empty
	// struct, the instantiated type alone determines its identity: values stored
	// with different type arguments never collide, and no exported key values are
	// needed.
	key[T any] struct{}

	// box wraps a stored value so the context always holds a non-nil interface,
	// even when T is an interface type and its value is nil. Get can therefore
	// distinguish a stored nil from a missing value.
	box[T any] struct{ data T }
)

// With returns a copy of ctx that carries value under a key derived from the
// type T.
//
// The key is the empty value key[T]{}, whose type is instantiated with T, so
// values of different types never collide, even when their underlying types are
// identical: with userID and teamID both defined as int, With[userID] and
// Get[teamID] stay independent. Each type holds at most one value in a context
// chain. Calling With again with the same T in a derived context overrides the
// value visible through that context, while values of other types are inherited
// from the parent.
//
// Any value of T is stored as-is, including zero values, nil pointers, and nil
// interface values, and can be retrieved with [Get]. Stored values are assigned
// with ordinary Go assignment and are not deep-copied: pointers, maps, and slice
// backing arrays may refer to data shared with the caller and with every context
// derived from the result. The package does not synchronize access to mutable
// data, so callers are responsible for safe concurrent access.
//
// The returned context is created with [context.WithValue], so it keeps the
// parent's deadline, cancellation, and value chain. With panics if ctx is nil.
//
// Because the key is derived from the static type T, retrieval must use the
// same type argument. Storing a value as an interface type, such as error or
// any, makes it retrievable only with that interface type, while type aliases
// share a key with their target type.
//
// Prefer passing data explicitly when the call chain allows it; use With for
// values that must cross API boundaries whose signatures cannot carry them,
// such as middleware and request handlers.
func With[T any](ctx context.Context, value T) context.Context {
	// Keep the stored interface non-nil even when value is a nil interface.
	// Get relies on this to distinguish a stored nil from an absent key.
	return context.WithValue(
		ctx,
		key[T]{}, box[T]{value},
	)
}

// Get returns the value of type T stored in ctx by [With] and reports whether
// it is present.
//
// The lookup matches only the exact type argument T used with With. Values
// stored under a different type argument are not returned, even when the types
// have the same underlying type or when one is an interface implemented by a
// value stored under another type.
//
// If no value of type T is present, Get returns the zero value of T and false.
// A stored value is always reported as present, even when it is a zero value, a
// nil pointer, or a nil interface value, so the two cases stay distinguishable.
// Values remain retrievable after the context has been cancelled.
//
// Get panics if ctx is nil.
func Get[T any](ctx context.Context) (T, bool) {
	stored, ok := ctx.Value(key[T]{}).(box[T])
	if !ok {
		var zero T
		return zero, false
	}
	return stored.data, true
}
