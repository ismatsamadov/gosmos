---
slug: your-first-check
title: Your First Check
type: code_tests
position: 2
difficulty: 5
topics: functions, testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Mirela liked your cases. Now she wants the checking code from you as well, starting with `batteryBand`, which every pre-pass poll relies on.

A checking function takes the test's `*testing.T`. When a result is wrong, `t.Errorf` marks the test failed and records a message, by convention `call = got, want expected`; then the test carries on:

```go
func checkCode(t *testing.T,
	code func(name string) string,
) {
	if got := code("Port Ellery"); got != "PEL" {
		t.Errorf("code(%q) = %q, want %q",
			"Port Ellery", got, "PEL")
	}
}
```

## Your task

Fill in `checkBand` in `checks.go`: call `band` at 8.2, 7.2 and 6.6 V and 0.01 V below each, and report each wrong band with `t.Errorf`. The tests run it against the real `batteryBand`, where it must stay quiet, and against three mutants, where it must complain.

> [!NOTE]
> `t.Fatalf` reports the same way but also stops the test, for when the checks after it would make no sense.
