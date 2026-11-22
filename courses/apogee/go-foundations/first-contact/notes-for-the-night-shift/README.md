---
slug: notes-for-the-night-shift
title: Notes for the Night Shift
type: code
position: 2
difficulty: 1
topics: comments
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Nkem reviews your first pass log: "Customers read this. Switch the DEBUG lines off, but keep them for the night shift."

Go ignores everything from `//` to the end of the line. Use comments for notes to people, or to switch a line off without deleting it:

```go
// Pipit-1 beacons every 30 seconds.
// fmt.Println("DEBUG: beacon timer started")
```

## Your task

Turn the two `DEBUG` lines into comments. The other three lines should still print, in the same order.

> [!NOTE]
> The moment the satellite leaves the station's view and the link ends is called LOS, short for loss of signal.
