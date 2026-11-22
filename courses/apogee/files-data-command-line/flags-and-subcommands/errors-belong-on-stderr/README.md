---
slug: errors-belong-on-stderr
title: Errors Belong on Stderr
type: code_tests
position: 6
difficulty: 6
topics: errors, programs, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Before calibration day, Mirela wants every gctl command to end the same way, through one `report` function. "A corrupted booking log must never look like a day with no passes. Show me the failure."

A script reads every line on standard output as a result. So a failed command writes no results, not even those found before the failure: only one line on standard error, and status 1. An empty report is still a success: no lines, status 0.

```go
var out, diag bytes.Buffer
args := []string{"gone.json"}
in := strings.NewReader("")
status := runPasses(args, in, &out, &diag)
fmt.Println(status, out.Len()) // 1 0
fmt.Print(diag.String())
// gctl: open gone.json: no such file or directory
```

## Your task

Implement `report`, which `runPasses` already ends with. When `err` is not nil, write only `gctl: ` and the error as one line to `stderr` and return 1, even if `lines` holds results. Otherwise write each line to `stdout` and return 0.
