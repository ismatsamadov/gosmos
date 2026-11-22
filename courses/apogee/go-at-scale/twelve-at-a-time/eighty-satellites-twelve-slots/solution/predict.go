package main

import (
	"context"
	"sync"
)

// predictAll predicts every satellite of els in a group of slots, and returns their passes by satellite ID, or the group's first failure, which stops the rest.
func predictAll(ctx context.Context, els []Elements, slots int, predict func(context.Context, Elements) ([]PassView, error)) (map[string][]PassView, error) {
	g, _ := newGroup(ctx, slots)
	var mu sync.Mutex
	byID := make(map[string][]PassView, len(els))
	for _, e := range els {
		g.Go(func(ctx context.Context) error {
			passes, err := predict(ctx, e)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			byID[e.SatID] = passes
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return byID, nil
}
