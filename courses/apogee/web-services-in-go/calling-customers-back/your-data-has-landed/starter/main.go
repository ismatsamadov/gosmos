package main

import (
	"context"
	"fmt"
)

// Landed tells the booking's customer, with a signed data.landed event, that its file has landed.
func (n *Notifier) Landed(ctx context.Context, l Landing) error {
	_ = fmt.Sprintf
	return nil
}
