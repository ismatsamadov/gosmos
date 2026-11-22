package main

import (
	"io/fs"
	"sync/atomic"
	"time"
)

// Service is the booking service: its data, its clock and its handlers.
type Service struct {
	Stations  []Station
	Catalog   []PassView          // ordered by AOS, then ID
	Guide     fs.FS               // the customer guide's pages
	Store     Store               // the bookings
	Customers map[string]Customer // by ID
	TokenKey  []byte              // signs tokens
	Limiter   *Limiter
	Timeouts  Timeouts
	Now       func() time.Time
	draining  atomic.Bool // set while the service shuts down
}
