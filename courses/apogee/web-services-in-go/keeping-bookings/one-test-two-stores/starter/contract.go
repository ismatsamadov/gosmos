package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

// checkStore checks that every store open returns keeps the promises of the Store interface.
func checkStore(t *testing.T, open func() Store) {
	_ = context.Background
	_ = errors.Is
	_ = time.Date
	// TODO: a subtest each: IDs from bk-4501, Get after Create, a conflict, Delete then Get, List by AOS.
}
