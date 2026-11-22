package main

import "slices"

// Customer is one company that books passes through the service.
type Customer struct {
	ID         string   // "furrowlight": lowercase, one word
	Name       string   // "Furrowlight"
	SatIDs     []string // the satellites it operates
	KeyHash    string   // lowercase hex SHA-256 of its whole API key
	WebhookURL string   // where its events go
	WebhookKey []byte   // the secret its events are signed with
}

// owns reports whether c operates the satellite satID.
func owns(c Customer, satID string) bool {
	return slices.Contains(c.SatIDs, satID)
}
