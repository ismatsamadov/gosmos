package main

import (
	"context"
	"sync"
)

// predictAll predicts every satellite of els, each while holding one of slots, and returns their passes by satellite ID, or the first error in els's order.
func predictAll(ctx context.Context, els []Elements, slots int, predict func(context.Context, Elements) ([]PassView, error)) (map[string][]PassView, error) {
	s := newSlots(slots)
	passes := make([][]PassView, len(els))
	errs := make([]error, len(els))
	var wg sync.WaitGroup
	for i, e := range els {
		if errs[i] = s.Acquire(ctx); errs[i] != nil {
			break // the run has ended: start no more
		}
		wg.Go(func() {
			defer s.Release()
			passes[i], errs[i] = predict(ctx, e)
		})
	}
	wg.Wait()
	byID := make(map[string][]PassView, len(els))
	for i, e := range els {
		if errs[i] != nil {
			return nil, errs[i]
		}
		byID[e.SatID] = passes[i]
	}
	return byID, nil
}
