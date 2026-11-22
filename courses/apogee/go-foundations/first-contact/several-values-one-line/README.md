---
slug: several-values-one-line
title: Several Values, One Line
type: code
position: 6
difficulty: 2
topics: printing
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
The short status message Pipit-1 broadcasts every 30 seconds, its beacon, is coming in at Port Ellery. The decoder prints every field on its own line, and Ingrid wants one line per fact.

`fmt.Println` takes any number of values, separated by commas. It prints them in order with a space between each, then ends the line. Text goes in quotes; numbers don't:

```go
fmt.Println("Ystmark wind", 64, "km/h")
```

That prints `Ystmark wind 64 km/h`. A satellite's mode is its operating state; `NOMINAL` means normal operation.

## Your task

Replace the seven calls in `main` with two, one for each line, passing the same values:

```text
Pipit-1 battery 7.91 V
Pipit-1 mode NOMINAL
```
