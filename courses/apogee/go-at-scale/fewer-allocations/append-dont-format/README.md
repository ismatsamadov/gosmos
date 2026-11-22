---
slug: append-dont-format
title: Append, Don't Format
type: code_tests
position: 3
difficulty: 9
topics: bytes, performance, time
created: 2026-11-12T03:10:33Z
modified: 2026-11-12T09:50:21Z
---
Tuesday afternoon. Each dish loader reads its day as text, one line a contact, and passplan 2 makes every line with `fmt.Sprintf`: six allocations a line, 12,120 for take-over day's plan.

A line holds the AOS, the LOS, the dish and the pass ID:

```text
05:29:26 05:37:52 1 furrow-a/VSK/24529
```

In February you appended pass IDs with `strconv.AppendInt`. Times and floats have appenders too: `t.AppendFormat` appends `t` to `dst` as `Format` would write it, and `strconv.AppendFloat` a float. While `dst` has room, neither allocates:

```go
dst := make([]byte, 0, 32)
dst = append(dst, "batt_v="...)
dst = strconv.AppendFloat(dst, 7.91, 'f', 2, 64)
dst = append(dst, " at "...)
dst = at.AppendFormat(dst, time.TimeOnly)
os.Stdout.Write(dst) // batt_v=7.91 at 05:25:30
```

## Your task

Make `appendContact` append each part of the line to `dst` itself, with no `fmt`: the AOS and LOS in the layout `time.TimeOnly` (`15:04:05`), the dish with `strconv.AppendInt`, the pass ID, the spaces and the newline. With room in `dst`, a line should allocate nothing.
