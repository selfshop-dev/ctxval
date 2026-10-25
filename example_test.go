package ctxval_test

import (
	"context"
	"fmt"

	"github.com/selfshop-dev/ctxval"
)

func Example() {
	type userID int
	type requestID string

	type tenant struct {
		ID   int
		Name string
	}

	// Build the request context.
	root := ctxval.With(context.Background(), userID(42))
	root = ctxval.With(root, requestID("req-001"))
	root = ctxval.With(root, tenant{ID: 7, Name: "acme"})

	// Create a derived context for a retry.
	// Only the requestID value is overridden.
	retry := ctxval.With(root, requestID("req-002"))

	rootUserID, _ := ctxval.Get[userID](root)
	rootRequestID, _ := ctxval.Get[requestID](root)

	retryRequestID, _ := ctxval.Get[requestID](retry)
	retryUserID, _ := ctxval.Get[userID](retry)
	retryTenant, _ := ctxval.Get[tenant](retry)

	_, hasLocale := ctxval.Get[string](retry)

	fmt.Printf("root user: %d\n", rootUserID)
	fmt.Printf("root request: %s\n", rootRequestID)
	fmt.Printf("retry user: %d\n", retryUserID)
	fmt.Printf("retry request: %s\n", retryRequestID)
	fmt.Printf("tenant: %s (%d)\n", retryTenant.Name, retryTenant.ID)
	fmt.Printf("locale present: %t\n", hasLocale)

	// Output:
	// root user: 42
	// root request: req-001
	// retry user: 42
	// retry request: req-002
	// tenant: acme (7)
	// locale present: false
}
