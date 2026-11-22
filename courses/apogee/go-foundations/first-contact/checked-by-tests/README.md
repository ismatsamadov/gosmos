---
slug: checked-by-tests
title: Checked by Tests
type: code_tests
position: 3
difficulty: 2
topics: functions, testing
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
The program that plans the day's passes, the scheduler, never reads what your program prints. It calls your functions, and a test says what it expects back.

This lesson comes with `main_test.go`, which you can read but not edit. A test calls a function, compares the result with what it wants, and reports a failure with `t.Errorf` when they differ:

```go
func TestBand(t *testing.T) {
	got := band()
	want := "UHF"
	if got != want {
		t.Errorf("band() = %q, want %q",
			got, want)
	}
}
```

Here, **Run** and **Submit** both run the tests.

## Your task

Open `main_test.go` and find what `TestStationCode` wants. Make `stationCode` in `main.go` return exactly that, by changing the text after `return`.
