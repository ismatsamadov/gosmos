package main

import (
	"reflect"
	"testing"
)

func TestPoll(t *testing.T) {
	var calls []int
	makeRule := func(id int, accept bool) Rule {
		return func(deg float64) bool {
			calls = append(calls, id)
			if deg != 41.5 {
				t.Errorf("rule argument = %v, want 41.5", deg)
			}
			return accept
		}
	}
	if got := poll(41.5, nil); !got {
		t.Error("poll(41.5, nil) = false, want true")
	}
	if got := poll(41.5, []Rule{makeRule(1, true), makeRule(2, false), makeRule(3, true)}); got {
		t.Error("poll(41.5, rejecting rules) = true, want false")
	}
	if !reflect.DeepEqual(calls, []int{1, 2}) {
		t.Errorf("rule calls = %v, want [1 2]", calls)
	}
	calls = nil
	if got := poll(41.5, []Rule{makeRule(1, true), makeRule(2, true)}); !got {
		t.Error("poll(41.5, accepting rules) = false, want true")
	}
	if !reflect.DeepEqual(calls, []int{1, 2}) {
		t.Errorf("rule calls = %v, want [1 2]", calls)
	}
}
