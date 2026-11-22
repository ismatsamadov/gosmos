---
slug: write-a-sheet-back
title: Write a Sheet Back
type: code_tests
position: 6
difficulty: 5
topics: errors, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Yusra needs the corrected booking sheet back to Hollis today, and his spreadsheet must read every note exactly as he typed it, commas, quotes and line breaks included.

A `csv.Writer` quotes each cell that needs it. It also buffers: `Write` may only fill memory, `Flush` sends what is left, and a failure on the way shows up in the writer's `Error` method, because `Flush` returns nothing.

```go
cw := csv.NewWriter(&out)
row := []string{"CWR", "STOWED, WIND 72"}
if err := cw.Write(row); err != nil {
	return err
}
cw.Flush()
return cw.Error() // out: CWR,"STOWED, WIND 72"
```

## Your task

Implement `writeSheet`: write every row to `w` with a `csv.Writer`, flush it, and return the writer's error, or the first error from `Write`.
