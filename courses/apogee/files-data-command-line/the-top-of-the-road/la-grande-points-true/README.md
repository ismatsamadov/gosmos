---
slug: la-grande-points-true
title: La Grande Points True
type: code_tests
position: 8
difficulty: 8
topics: errors, files, functions, programs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:19:20Z
---
Thursday, 06:00 on the salt flat, 13:00 UTC. Paz plugs in the stick with the settings, Yusra's export and Hollis's sheet: calibration day runs from these files and gctl alone.

Order the last steps so that a failure leaves nothing half done. `planDay` reads and checks everything first; then the report is saved whole; only a saved report is printed. Every path starts at `dir` and takes the day from the checked plan, never from the current folder:

```go
note := func(w io.Writer) error {
	return writeNote(w, stow)
}
if err := keepNote(path, note); err != nil {
	return err // a note not kept is not shown
}
return writeNote(stdout, stow)
```

## Your task

Complete `runDay`:

1. Plan the day with `planDay`.
2. Make `reports/daily` below `dir` with mode `0700`, and save the report there as `<day>.txt` with `saveDay` and `writeReport`.
3. Only then write the report to `stdout`. Return the first error.
