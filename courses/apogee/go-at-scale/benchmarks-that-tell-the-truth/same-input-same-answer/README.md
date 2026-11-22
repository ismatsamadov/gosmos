---
slug: same-input-same-answer
title: Same Input, Same Answer
type: code_tests
position: 5
difficulty: 8
topics: errors, performance, slices
created: 2026-11-11T22:45:16Z
modified: 2026-11-11T22:45:16Z
---
Wednesday 16 May, 13:30. Brinewatch's thirty-six satellites leave their rocket, and before anyone times the two planners on the days ahead, Nkem wants to know that they plan the same contacts.

A version that gives a different answer isn't faster, it's wrong. So before timing two versions, run both on the same input and check that they agree, saying where they don't:

```go
want := sortSlow(slices.Clone(gusts))
got := sortFast(slices.Clone(gusts))
if !slices.Equal(got, want) {
	b.Fatal("sortFast disagrees with sortSlow")
}
// only now, for b.Loop()
```

## Your task

`sameContacts` returns `nil` for any two lists. When the lengths of `old` and `new` differ, return an error that says so; otherwise return one for the first index whose pass ID or dish differs, giving the index and each list's pass ID and dish, or `nil` when none differs. Use these two forms:

```text
old has %d contacts, new has %d
contact %d: old %s dish %d, new %s dish %d
```

> [!NOTE]
> The tests take Thursday 17 May as `planOld` and `planNew` both plan it, 1,231 contacts, and plant differences in a copy.
