---
slug: stow-a-station
title: Stow a Station
type: code_tests
position: 3
difficulty: 6
topics: errors, pointers, slices
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T17:14:45Z
---
Nobody lives at Cape Wren. When the wind reaches 72 km/h, its dish parks itself pointing up, which ops calls a stow. Ingrid wants each stow recorded in gctl, with its reason, so later commands know.

Ranging over a slice gives you a copy of each element. To change an element in place, go through its index:

```go
for i := range readings {
	readings[i].Checked = true
}
```

## Your task

Write `cmdStow`: mark the station whose code follows `stow` as stowed, keep the words after the code as its reason, and print its stow note. A missing code, an unknown one and a station already stowed are errors that leave the table as it was. Add `stow` to the table. Among the lines gctl then prints:

```text
CWR stowed: wind 74 km/h
error: stow needs a station code
error: unknown station "XYZ"
error: CWR is already stowed
stow      record a stow: stow CODE [reason...]
```
