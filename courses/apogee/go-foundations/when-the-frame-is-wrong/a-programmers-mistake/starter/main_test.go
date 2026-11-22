package main

import "testing"

func TestRequireLimit(t *testing.T) {
	for _, limit := range []int{-1, 0, 1, 12} {
		func() {
			defer func() {
				got := recover()
				if limit <= 0 {
					if got != "invalid decoder limit" {
						t.Errorf("requireLimit(%d) panic = %v, want invalid decoder limit", limit, got)
					}
				} else if got != nil {
					t.Errorf("requireLimit(%d) panic = %v, want nil", limit, got)
				}
			}()
			requireLimit(limit)
		}()
	}
}
