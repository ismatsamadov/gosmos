---
slug: the-frame-validator
title: The Frame Validator
type: code_tests
position: 12
difficulty: 7
topics: errors, slices
created: 2026-11-01T07:36:46Z
modified: 2026-11-04T17:00:51Z
---
Thistle's frame 4182 arrived twice, and the second copy also claims 900 °C. The validator stops at the first problem, so the log says only `duplicate frame`.

`errors.Join` combines several errors into one. Its message puts each on its own line, and `errors.Is` and `errors.AsType` find every one inside. Given no errors, it returns `nil`:

```go
err := errors.Join(
	errors.New("radome heater off"),
	errors.New("fiber down"),
)
fmt.Println(err) // two lines, one per error
```

## Your task

Make `validate` report every problem: collect each one in a slice instead of returning it, then return them joined. A healthy frame must still give `nil`.
