---
slug: errors-are-values
title: Errors Are Values
type: code_tests
position: 0
difficulty: 4
topics: errors
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Thistle's battery read 6.2 V last night, and the telemetry service stored it without a word. Ingrid's new rule: every failure gets a name, a reason and a line in the log.

A Go function that can fail returns an `error`: an interface with one method, `Error() string`. `nil` means all went well; anything else describes the failure. `errors.New` makes one from a message:

```go
func checkWind(kmh float64) error {
	if kmh >= 72 {
		return errors.New("wind too strong")
	}
	return nil
}
```

## Your task

Make `checkBattery` return an error reading `battery too low` for readings under `minBattV`, and `nil` otherwise. `errors.New` comes from the `errors` package, so add it to the imports.
