---
slug: quick-check-price-of-reflection
title: Quick Check: The Price of Reflection
type: quiz
position: 8
difficulty: 8
topics: methods, performance, reflection
created: 2026-11-12T03:14:43Z
modified: 2026-11-12T03:14:43Z
---
Thursday 7 June. Before the validator ships, Nkem asks why each settings struct shouldn't check itself with a method, and starts one on the whiteboard.

```go
func (h HostConfig) Validate() error {
	if h.Dishes < 1 || h.Dishes > 2 {
		return fmt.Errorf("dishes: %d", h.Dishes)
	}
	return nil // and six more fields
}
```

She has timed a draft of the validator, which takes about 1 µs and 38 allocations to check one `Config`; a method written for `Config` takes 2 ns and none.

## Your task

Answer her three questions: when reading tags pays, when a tag's mistake shows up, and where the check belongs.
