package ctxval_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/selfshop-dev/ctxval"
)

func assertRoundTrip[T any](t *testing.T, want T) {
	t.Helper()

	ctx := ctxval.With(context.Background(), want)

	got, ok := ctxval.Get[T](ctx)

	require.True(t, ok, "value should exist in context")
	assert.Equal(t, want, got)
}

func TestWithAndGetRoundTrip(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assertRoundTrip(t, 42)
	})

	t.Run("string", func(t *testing.T) {
		assertRoundTrip(t, "hello")
	})

	t.Run("struct", func(t *testing.T) {
		type user struct {
			ID   int
			Name string
		}

		assertRoundTrip(t, user{
			ID:   42,
			Name: "Alice",
		})
	})

	t.Run("slice", func(t *testing.T) {
		assertRoundTrip(t, []string{"read", "write"})
	})
}

func TestGetMissingReturnsZeroAndFalse(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		got, ok := ctxval.Get[int](context.Background())

		assert.False(t, ok)
		assert.Zero(t, got)
	})

	t.Run("struct", func(t *testing.T) {
		type state struct {
			Name    string
			Enabled bool
		}

		got, ok := ctxval.Get[state](context.Background())

		assert.False(t, ok)
		assert.Equal(t, state{}, got)
	})
}

func TestStoredZeroValuesArePresent(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		ctx := ctxval.With(context.Background(), 0)

		got, ok := ctxval.Get[int](ctx)

		require.True(t, ok)
		assert.Equal(t, 0, got)
	})

	t.Run("string", func(t *testing.T) {
		ctx := ctxval.With(context.Background(), "")

		got, ok := ctxval.Get[string](ctx)

		require.True(t, ok)
		assert.Empty(t, got)
	})

	t.Run("bool", func(t *testing.T) {
		ctx := ctxval.With(context.Background(), false)

		got, ok := ctxval.Get[bool](ctx)

		require.True(t, ok)
		assert.False(t, got)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var value *int
		ctx := ctxval.With(context.Background(), value)

		got, ok := ctxval.Get[*int](ctx)

		require.True(t, ok)
		assert.Nil(t, got)
	})
}

func TestDifferentTypesDoNotCollide(t *testing.T) {
	type userID int
	type teamID int

	ctx := ctxval.With(context.Background(), userID(42))

	// The underlying types are both int, but the named types differ.
	gotTeamID, ok := ctxval.Get[teamID](ctx)

	assert.False(t, ok)
	assert.Equal(t, teamID(0), gotTeamID)

	ctx = ctxval.With(ctx, teamID(7))

	gotUserID, userOK := ctxval.Get[userID](ctx)
	gotTeamID, teamOK := ctxval.Get[teamID](ctx)

	require.True(t, userOK)
	require.True(t, teamOK)
	assert.Equal(t, userID(42), gotUserID)
	assert.Equal(t, teamID(7), gotTeamID)
}

func TestWithOverridesOnlyTheSameTypeInChildContext(t *testing.T) {
	type traceID string
	type userID int

	parent := ctxval.With(context.Background(), traceID("trace-parent"))
	parent = ctxval.With(parent, userID(42))

	child := ctxval.With(parent, traceID("trace-child"))

	parentTrace, ok := ctxval.Get[traceID](parent)
	require.True(t, ok)
	assert.Equal(t, traceID("trace-parent"), parentTrace)

	childTrace, ok := ctxval.Get[traceID](child)
	require.True(t, ok)
	assert.Equal(t, traceID("trace-child"), childTrace)

	// The value of another type is inherited by the child.
	childUserID, ok := ctxval.Get[userID](child)
	require.True(t, ok)
	assert.Equal(t, userID(42), childUserID)
}

func TestWithPreservesParentContextCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx := ctxval.With(parent, "value")

	cancel()

	require.ErrorIs(t, ctx.Err(), context.Canceled)

	select {
	case <-ctx.Done():
		// The derived context propagates cancellation.
	default:
		t.Fatal("expected ctx.Done() to be closed")
	}

	got, ok := ctxval.Get[string](ctx)
	require.True(t, ok)
	assert.Equal(t, "value", got)
}

func TestWithPreservesParentDeadline(t *testing.T) {
	deadline := time.Now().Add(time.Hour)

	parent, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	ctx := ctxval.With(parent, "value")

	got, ok := ctx.Deadline()
	require.True(t, ok)
	assert.Equal(t, deadline, got)
}

func TestUnrelatedContextValuesAreInherited(t *testing.T) {
	type parentKey struct{}

	parent := context.WithValue(context.Background(), parentKey{}, "parent")
	ctx := ctxval.With(parent, 42)

	assert.Equal(t, "parent", ctx.Value(parentKey{}))

	got, ok := ctxval.Get[int](ctx)
	require.True(t, ok)
	assert.Equal(t, 42, got)
}

func TestStoredNilInterfaceIsPresent(t *testing.T) {
	var value error

	ctx := ctxval.With(context.Background(), value)

	got, ok := ctxval.Get[error](ctx)

	require.True(t, ok)
	assert.NoError(t, got)
}

func TestInterfaceAndConcreteValuesUseDifferentSlots(t *testing.T) {
	t.Run("stored as interface", func(t *testing.T) {
		ctx := ctxval.With[any](context.Background(), "value")

		got, ok := ctxval.Get[any](ctx)
		require.True(t, ok)
		assert.Equal(t, "value", got)

		_, ok = ctxval.Get[string](ctx)
		assert.False(t, ok)
	})

	t.Run("stored as concrete type", func(t *testing.T) {
		ctx := ctxval.With(context.Background(), "value")

		_, ok := ctxval.Get[any](ctx)
		assert.False(t, ok)
	})
}

func TestTypeAliasSharesSlotWithTarget(t *testing.T) {
	type alias = int

	ctx := ctxval.With(context.Background(), 42)

	got, ok := ctxval.Get[alias](ctx)

	require.True(t, ok)
	assert.Equal(t, 42, got)
}

func TestNilContextPanics(t *testing.T) {
	t.Run("With", func(t *testing.T) {
		// A typed nil interface value is passed deliberately: the call must panic
		// because the contract requires a non-nil context.
		var ctx context.Context

		assert.Panics(t, func() {
			ctxval.With(ctx, 1)
		})
	})

	t.Run("Get", func(t *testing.T) {
		var ctx context.Context

		assert.Panics(t, func() {
			ctxval.Get[int](ctx)
		})
	})
}

func TestConcurrentReads(t *testing.T) {
	ctx := ctxval.With(context.Background(), 42)

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			got, ok := ctxval.Get[int](ctx)
			assert.True(t, ok)
			assert.Equal(t, 42, got)
		}()
	}
	wg.Wait()
}
