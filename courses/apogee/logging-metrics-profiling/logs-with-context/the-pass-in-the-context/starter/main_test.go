package main

import (
	"context"
	"testing"
	"time"
)

func TestPassInContext(t *testing.T) {
	const cirrine = "cirrine-1/LAN/16402"
	bg := context.Background()
	t.Run("there and back", func(t *testing.T) {
		if id, ok := passFrom(withPass(bg, cirrine)); id != cirrine || !ok {
			t.Errorf("passFrom(withPass(ctx, %q)) = %q, %v, want %q, true", cirrine, id, ok, cirrine)
		}
	})
	t.Run("a child context still finds it", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(withPass(bg, cirrine), time.Minute)
		defer cancel()
		if id, ok := passFrom(ctx); id != cirrine || !ok {
			t.Errorf("passFrom(a child of withPass(ctx, %q)) = %q, %v, want %q, true", cirrine, id, ok, cirrine)
		}
	})
	t.Run("no pass, no ID", func(t *testing.T) {
		if id, ok := passFrom(bg); id != "" || ok {
			t.Errorf(`passFrom(context.Background()) = %q, %v, want "", false`, id, ok)
		}
	})
	t.Run("an empty ID is none", func(t *testing.T) {
		if id, ok := passFrom(withPass(bg, "")); id != "" || ok {
			t.Errorf(`passFrom(withPass(ctx, "")) = %q, %v, want "", false`, id, ok)
		}
	})
	t.Run(`a value under the string "pass_id" is not yours`, func(t *testing.T) {
		ctx := context.WithValue(bg, "pass_id", cirrine)
		if id, ok := passFrom(ctx); id != "" || ok {
			t.Errorf(`passFrom(a context holding "pass_id" as a plain string key) = %q, %v, want "", false`, id, ok)
		}
	})
}
