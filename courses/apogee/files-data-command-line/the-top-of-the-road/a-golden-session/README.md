---
slug: a-golden-session
title: A Golden Session
type: code_tests
position: 6
difficulty: 7
topics: programs, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:19:20Z
---
Ingrid wants to know, byte for byte, what gctl prints on Thursday. "Copy. Exact output, checked before the first AOS."

A golden test runs a whole session and compares all of it, standard output, standard error and the exit status, with a file kept in `testdata`. The golden file is the specification: open one to see exactly what a session must print. One extra newline fails the test:

```go
got := transcript(stdout.String(),
	stderr.String(), code)
file := filepath.Join("testdata", name+".golden")
want, err := os.ReadFile(file)
```

## Your task

Implement `run` so every session matches its golden file. `parseArgs` returns the options, the usage that fits the arguments (gctl's own or the day command's) and an error:

1. For `flag.ErrHelp`, print that usage to standard output and return 0.
2. For any other error, print the error on a line of its own and then that usage to standard error, and return 2.
3. Otherwise plan the day with `planDay` and write it with `writeReport`. If either fails, print `gctl: <error>` to standard error and return 1.
