---
slug: pre-pass-checklist
title: The Pre-Pass Checklist
type: code
position: 1
difficulty: 1
topics: printing, statements
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Pipit-1 rises over Port Ellery in ten minutes, and the checklist printer has shuffled Ingrid's steps. The moment the satellite comes into view is called AOS, short for acquisition of signal.

Go runs the statements in `main` from top to bottom, one at a time. Each `fmt.Println` prints one line, so the order of the calls is the order of the lines:

```go
fmt.Println("Antenna parked")
fmt.Println("Antenna moving to the rise azimuth")
```

## Your task

Reorder the `fmt.Println` calls so the title prints first, then steps 1 to 4 in order. Move whole lines; don't change their text.

