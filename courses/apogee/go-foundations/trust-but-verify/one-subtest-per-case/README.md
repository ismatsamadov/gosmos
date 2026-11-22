---
slug: one-subtest-per-case
title: One Subtest per Case
type: code_tests
position: 4
difficulty: 6
topics: testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-22T10:03:44Z
---
One of the Cape Wren wind cases failed after a refactor, a change to the code's structure without changing its intended behavior. The log named only `TestWind`, not the case. Nkem asks which one it was.

`t.Run(name, f)` runs `f` as a subtest with its own name and its own `*testing.T`. A failure inside it is reported as `TestWind/<name>`, and the other subtests still run:

```go
for _, code := range []string{"YST", "CWR"} {
	t.Run(code, func(t *testing.T) {
		if !known(code) {
			t.Error("known = false, want true")
		}
	})
}
```

## Your task

In `checkWind`, run each case as a subtest named `c.name`, with the check inside it. Leave the cases as they are.

> [!TIP]
> Inside the subtest, report with the `t` that `t.Run` hands you, not the outer one.
