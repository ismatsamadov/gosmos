package main

import "context"

// Customer is one company that books passes through the service.
type Customer struct {
	ID         string   // "furrowlight": lowercase, one word
	Name       string   // "Furrowlight"
	SatIDs     []string // the satellites it operates
	KeyHash    string   // lowercase hex SHA-256 of its whole API key
	WebhookURL string   // where its events go
	WebhookKey []byte   // the secret its events are signed with
}

// customerKey is the type of the context key a request's customer is kept under.
type customerKey struct{}

// withCustomer returns a copy of ctx that carries c.
func withCustomer(ctx context.Context, c Customer) context.Context {
	return context.WithValue(ctx, customerKey{}, c)
}

// customerFrom returns the customer ctx carries, if authenticate put one there.
func customerFrom(ctx context.Context) (Customer, bool) {
	c, ok := ctx.Value(customerKey{}).(Customer)
	return c, ok
}
