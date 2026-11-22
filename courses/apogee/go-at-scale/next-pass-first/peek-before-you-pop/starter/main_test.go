package main

import (
	"cmp"
	"fmt"
	"testing"
)

// safePeek peeks at q, turning a panic into an error, so a test can report it.
func safePeek[T any](q *Queue[T]) (v T, ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	v, ok = q.Peek()
	return v, ok, nil
}

func TestPeek(t *testing.T) {
	t.Run("minutes to the next pass of three streams: 5, 2 and 8", func(t *testing.T) {
		q := NewQueue(cmp.Compare[int])
		for _, v := range []int{5, 2, 8} {
			q.Push(v)
		}
		for i := range 2 {
			if got, ok := q.Peek(); got != 2 || !ok {
				t.Fatalf("Peek() number %d = %d, %v, want 2, true", i+1, got, ok)
			}
		}
		if q.Len() != 3 {
			t.Errorf("Len() after two calls to Peek = %d, want 3: a peek takes nothing out", q.Len())
		}
	})
	t.Run("after a Pop", func(t *testing.T) {
		q := NewQueue(cmp.Compare[int])
		for _, v := range []int{5, 2, 8} {
			q.Push(v)
		}
		q.Pop()
		if got, ok := q.Peek(); got != 5 || !ok {
			t.Errorf("Peek() after popping 2 = %d, %v, want 5, true", got, ok)
		}
	})
	t.Run("an empty queue of ints", func(t *testing.T) {
		got, ok, err := safePeek(NewQueue(cmp.Compare[int]))
		if err != nil {
			t.Fatalf("Peek() on an empty queue panicked: %v", err)
		}
		if got != 0 || ok {
			t.Errorf("Peek() on an empty queue of ints = %d, %v, want 0, false", got, ok)
		}
	})
	t.Run("an empty queue of station codes", func(t *testing.T) {
		got, ok, err := safePeek(NewQueue(cmp.Compare[string]))
		if err != nil {
			t.Fatalf("Peek() on an empty queue panicked: %v", err)
		}
		if got != "" || ok {
			t.Errorf("Peek() on an empty queue of strings = %q, %v, want \"\", false", got, ok)
		}
	})
}
