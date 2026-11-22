---
slug: known-failures
title: Known Failures
type: code_tests
position: 4
difficulty: 5
topics: errors
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Thistle's frames often arrive twice, once through each of two stations. Duplicates are normal, so the pager should stay quiet about them and page for anything else.

A *sentinel error* is one error value, declared once, that callers recognize. Return that value itself: two `errors.New` calls with the same text still make two different errors. `errors.Is(err, target)` reports whether `err` is the target:

```go
var ErrStowed = errors.New("dish stowed")

if errors.Is(err, ErrStowed) {
	fmt.Println("waiting for the wind to drop")
}
```

## Your task

Make `accept` return `ErrDuplicate` itself for a repeated frame. Then make `shouldPage` return `false` for `nil` and for `ErrDuplicate`, and `true` for any other error.

> [!NOTE]
> Why `errors.Is` and not `==`? Next lesson.
