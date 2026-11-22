package main

import (
	"context"
	"sync"
)

// Group runs jobs in goroutines, at most limit at once, and stops at its first failure.
type Group struct {
	ctx    context.Context
	cancel context.CancelCauseFunc
	slots  chan struct{} // one value for each job running
	wg     sync.WaitGroup
	once   sync.Once
	err    error // the first failure
}

// newGroup returns a group that runs at most limit jobs at once, and the context its jobs get, which ends when ctx does or at the group's first failure.
func newGroup(ctx context.Context, limit int) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{ctx: ctx, cancel: cancel, slots: make(chan struct{}, limit)}, ctx
}

// Go waits for a free slot and runs job in a goroutine of its own; once the group's context has ended, it runs nothing and fails the group with the cause.
func (g *Group) Go(job func(ctx context.Context) error) {
	g.slots <- struct{}{}
	g.wg.Go(func() {
		defer func() { <-g.slots }()
		job(g.ctx) // TODO: fail the group when job returns an error
	})
}

// Wait waits for every job to return, then returns the group's first failure.
func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel(g.err) // every job is done: let the group's context go
	return g.err
}
