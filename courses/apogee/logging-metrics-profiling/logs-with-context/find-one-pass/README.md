---
slug: find-one-pass
title: Find One Pass
type: code_tests
position: 5
difficulty: 7
topics: json, logging, streams
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Thursday afternoon. Next time a pass fails, the night shift must find its lines in a minute, not a Monday: one search through every program's log for one pass ID.

`json.Unmarshal` into a struct fills the fields the struct names and ignores every other key, so `Line`, in `line.go`, reads the six keys every program writes from any line. A `bufio.Scanner` hands you the input one line at a time; a line that doesn't decode is counted and skipped, never a reason to stop:

```go
var l Line
err := json.Unmarshal([]byte(
	`{"msg":"pass started","seq":4}`), &l)
fmt.Printf("%q %v\n", l.Msg, err)
// "pass started" <nil>
```

## Your task

Make `follow` read `r` line by line and return, in the order read, the lines whose `pass_id` is `passID`, each with `Raw` set to the line as written. Count the lines that aren't JSON in `broken`, and return the scanner's error, if any.

> [!NOTE]
> The test reads `testdata/bench.log`: the four programs' lines from a first run of Nkem's bench replay of 8 November, one log after another.
