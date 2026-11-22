package main

import (
	"context"
	"fmt"
	"time"
)

// passContext returns a context for one pass, which ends at the pass's LOS.
func passContext(parent context.Context, los time.Time) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

func main() {
	los := time.Date(2027, 6, 29, 4, 21, 0, 0, time.UTC)
	ctx, cancel := passContext(context.Background(), los)
	defer cancel()
	deadline, ok := ctx.Deadline()
	fmt.Println("pass ends at", deadline.Format(time.TimeOnly), ok)
}
