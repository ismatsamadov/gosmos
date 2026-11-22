---
slug: one-group-per-pass
title: One Group per Pass
type: code_tests
position: 4
difficulty: 7
topics: json, logging, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Wednesday, Nkem: "Every program describes a pass the same way, or no search finds them all." One `pass` object: its satellite, its station, and the moments it rises into view and sets, which ops calls AOS and LOS.

`slog.Group` gathers attributes under one key. The JSON handler writes a group as a nested object, and the text handler writes each of its attributes with the group's name and a dot in front:

```go
g := slog.Group("radio",
	slog.String("band", "S"),
	slog.Float64("margin_db", -1.5))
l.Warn("carrier lost", g)
// JSON: "radio":{"band":"S","margin_db":-1.5}
// text: radio.band=S radio.margin_db=-1.5
```

## Your task

Make `passGroup` return the group `pass` with four attributes, in this order: `sat_id` and `station` as strings, then `aos` and `los` as times in UTC. Cirrine-1's pass of 8 November reads, on one line:

```text
"pass":{"sat_id":"cirrine-1","station":"LAN",
"aos":"2028-11-08T01:55:09Z",
"los":"2028-11-08T02:02:39Z"}
```

> [!TIP]
> `Time.UTC` gives the same moment in UTC, so a time read with another zone's offset still prints as the Board shows it.
