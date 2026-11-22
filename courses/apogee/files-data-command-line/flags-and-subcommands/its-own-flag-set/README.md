---
slug: its-own-flag-set
title: Its Own Flag Set
type: code_tests
position: 1
difficulty: 6
topics: errors, programs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
In the van, Paz only needs Salar Quieto's bookings, and her old script took the station as its first word. In gctl it becomes a flag, an option named with a dash: `-station SQT`.

The `flag` package reads flags. Each command makes its own `flag.NewFlagSet`, with a name and `flag.ContinueOnError`, so a mistake comes back from `Parse` as an error instead of ending the program. `SetOutput` chooses where the flag set writes its messages.

```go
fs := flag.NewFlagSet("sweep",
	flag.ContinueOnError)
fs.SetOutput(&messages)
step := fs.Float64("step", 0.5, "step in degrees")
err := fs.Parse([]string{"-step", "0.25"})
fmt.Println(*step, err) // 0.25 <nil>
```

## Your task

Implement `parseOptions`: make a flag set named `passes` that continues on error and writes its messages to `output`, define `-station` with an empty default, meaning every station, then parse `args`. Return the options, or empty options and the parse error.
