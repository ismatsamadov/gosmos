---
slug: first-verbs
title: First Verbs
type: code
position: 0
difficulty: 2
topics: formatting, printing
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Ingrid wants the Board, the big screen in Apogee's glass-walled operations room, the Fishbowl, to open with one sentence. Right now its header comes out in three pieces, one per line.

`fmt.Printf` prints a format string. Each verb in it, a `%` and a letter, is replaced by the next value you pass: `%s` takes a string, `%d` a whole number. `Printf` adds no newline, so end the format with `\n`:

```go
station := "Cape Wren"
windKmh := 48
fmt.Printf("%s wind: %d km/h\n", station, windKmh)
```

That prints `Cape Wren wind: 48 km/h`.

## Your task

Replace the first three `Println` calls in `main` with one `Printf` that uses `name` and `passes`, and keep the last call. The Board then reads:

```text
Next passes for Pipit-1: 4
First AOS 04:12 UTC over Port Ellery
```
