package ctxval_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/selfshop-dev/ctxval"
)

func TestWithAndGet(t *testing.T) {
	ctx := context.Background()
	ctx = ctxval.With(ctx, "hello")

	got, ok := ctxval.Get[string](ctx)
	require.True(t, ok)
	assert.Equal(t, "hello", got)
}

func TestGetMissing(t *testing.T) {
	ctx := context.Background()

	got, ok := ctxval.Get[int](ctx)
	assert.False(t, ok)
	assert.Zero(t, got)
}

func TestGetOr(t *testing.T) {
	t.Run("absent returns fallback", func(t *testing.T) {
		ctx := context.Background()
		got := ctxval.GetOr(ctx, 99)
		assert.Equal(t, 99, got)
	})

	t.Run("present returns stored value", func(t *testing.T) {
		ctx := ctxval.With(context.Background(), 42)
		got := ctxval.GetOr(ctx, 99)
		assert.Equal(t, 42, got)
	})
}

func TestExists(t *testing.T) {
	ctx := context.Background()
	assert.False(t, ctxval.Exists[string](ctx))

	ctx = ctxval.With(ctx, "present")
	assert.True(t, ctxval.Exists[string](ctx))
	assert.False(t, ctxval.Exists[int](ctx))
}

func TestMustPresent(t *testing.T) {
	ctx := ctxval.With(context.Background(), time.Second)
	got := ctxval.Must[time.Duration](ctx)
	assert.Equal(t, time.Second, got)
}

func TestMustPanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = ctxval.Must[string](context.Background())
	})
}

func TestMustPanicMessage(t *testing.T) {
	assert.PanicsWithValue(t,
		"ctxval.Must: value of type string not found in context",
		func() {
			_ = ctxval.Must[string](context.Background())
		},
	)
}

func TestShadowing(t *testing.T) {
	ctx := context.Background()
	ctx = ctxval.With(ctx, "first")
	ctx = ctxval.With(ctx, "second")

	got, ok := ctxval.Get[string](ctx)
	require.True(t, ok)
	assert.Equal(t, "second", got)
}

func TestDifferentTypesIndependent(t *testing.T) {
	ctx := context.Background()
	ctx = ctxval.With(ctx, "str")
	ctx = ctxval.With(ctx, 7)
	ctx = ctxval.With(ctx, true)

	s, ok := ctxval.Get[string](ctx)
	require.True(t, ok)
	assert.Equal(t, "str", s)

	n, ok := ctxval.Get[int](ctx)
	require.True(t, ok)
	assert.Equal(t, 7, n)

	b, ok := ctxval.Get[bool](ctx)
	require.True(t, ok)
	assert.True(t, b)
}

func TestTypedNil(t *testing.T) {
	type User struct{ Name string }

	var u *User // typed nil
	ctx := ctxval.With(context.Background(), u)

	got, ok := ctxval.Get[*User](ctx)
	require.True(t, ok, "typed nil must be reported as present")
	assert.Nil(t, got)
}

func TestParentChain(t *testing.T) {
	parent := ctxval.With(context.Background(), "parent-value")
	child := ctxval.With(parent, 100)

	// Child still sees the parent's string.
	s, ok := ctxval.Get[string](child)
	require.True(t, ok)
	assert.Equal(t, "parent-value", s)

	// Parent does not see the child's int.
	_, ok = ctxval.Get[int](parent)
	assert.False(t, ok)
}

func TestShadowingPreservesParent(t *testing.T) {
	parent := ctxval.With(context.Background(), "parent")
	child := ctxval.With(parent, "child")

	// Child sees the shadowed value.
	got, ok := ctxval.Get[string](child)
	require.True(t, ok)
	assert.Equal(t, "child", got)

	// Parent still sees its original value.
	got, ok = ctxval.Get[string](parent)
	require.True(t, ok)
	assert.Equal(t, "parent", got)
}
