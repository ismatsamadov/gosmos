---
slug: save-only-a-complete-report
title: Save Only a Complete Report
type: code_tests
position: 4
difficulty: 7
topics: errors, files, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T17:07:41Z
---
Paz saved Thursday's plan last night. If this morning's save stops at its third line, the stick must still hold last night's plan, complete.

The report comes from a function that writes to an `io.Writer`, so it can fail after some lines are already out. Hand it a temporary file made beside the destination, never the destination itself, and rename only after the function and `Close` have both succeeded. On any failure, remove the temporary file, as you did for the schedule in July:

```go
var preview bytes.Buffer
if err := report(&preview); err != nil {
	return err // only part of the report
}
```

## Your task

Implement `saveDay(path, report)`: create a temporary file in the destination's folder with `os.CreateTemp`, call `report` with it, close it, then rename it over the destination. Return the first error, so a failed report leaves the destination as it was.
