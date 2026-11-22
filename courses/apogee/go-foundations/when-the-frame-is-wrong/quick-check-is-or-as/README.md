---
slug: quick-check-is-or-as
title: Quick Check: Is or As?
type: quiz
position: 7
difficulty: 5
topics: errors
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Mirela is reviewing the pager before it goes live, and she asks you to pick its tools. Both questions are about errors that reach the pager wrapped, like these two, where `checkTemp(51)` returns a `*LimitError`:

```go
dup := fmt.Errorf("thistle seq 4182: %w",
	ErrDuplicate)
hot := fmt.Errorf("thistle seq 4190: %w",
	checkTemp(51))
```

Choose the tool that keeps working when someone adds another layer of wrapping later.
