---
slug: seeds-that-catch-bugs
title: Seeds That Catch Bugs
type: code_tests
position: 1
difficulty: 8
topics: testing
created: 2026-11-11T12:51:44Z
modified: 2026-11-11T12:51:44Z
---
Mirela has planted one bug in each of three copies of her reference parser, the strict one she trusts: Saturday's two bugs and one more. She wants seeds that would have caught Saturday before Saturday.

Plain `go test` runs a fuzz test on its seed corpus alone: the values `f.Add` adds and the files under `testdata`. Invented inputs come only with `-fuzz`. So a broken parser fails `go test` only when a seed walks into its bug, and good inputs alone catch nothing. This test seeds two bad codes beside a good one:

```go
func FuzzStationCode(f *testing.F) {
	for _, s := range []string{"PEL", "", "pel"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		checkCode(t, s)
	})
}
```

## Your task

`lineSeeds`, in `seeds.go`, holds one good line, and `FuzzParseLine` adds every seed. Add a seed that makes it fail for each of Mirela's mutants: one panics on a line of five fields, one accepts `batt_v=nan`, one accepts `mode=nominal`.

> [!TIP]
> Copy the good line and change one thing in each copy: a line with two faults is refused for the first, and the second is never reached.
