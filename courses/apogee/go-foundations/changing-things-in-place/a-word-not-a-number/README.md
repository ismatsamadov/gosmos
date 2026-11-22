---
slug: a-word-not-a-number
title: A Word, Not a Number
type: code
position: 8
difficulty: 5
topics: methods, printing
created: 2026-11-01T07:36:46Z
modified: 2026-11-04T17:07:34Z
---
An alert at 03:00 says "severity 2", and the engineer on call, the person responsible for responding, has to look up what 2 means.

When a type has a method `String() string`, `Println`, `%v` and `%s` print what that method returns:

```go
type Wind float64

func (w Wind) String() string {
	return fmt.Sprintf("%.0f km/h", float64(w))
}

fmt.Println(Wind(81)) // 81 km/h
```

## Your task

Where the `TODO` is, add a `String` method to `Severity` that returns `INFO`, `WARNING` or `CRITICAL`, so the program prints:

```text
INFO
WARNING
CRITICAL
```

> [!WARNING]
> Give `String` a value receiver, as `Wind` does. `fmt.Println(Info)` hands `fmt` a copy of the value, not its address, so a `String` with a pointer receiver is never called.
