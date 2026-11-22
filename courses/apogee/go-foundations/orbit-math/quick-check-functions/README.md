---
slug: quick-check-functions
title: Quick Check: Functions
type: quiz
position: 4
difficulty: 3
topics: functions, multiple-returns
created: 2026-11-01T07:36:45Z
modified: 2026-11-22T10:08:26Z
---
Nkem is about to review the tool that sends commands from the ground to a satellite, the uplink. Its sequence counter must grow by one with every command. Before her review she asks three questions; the first is about this code:

```go
func nextSeq(seq int) int {
	seq = seq + 1
	return seq
}

seq := 117
nextSeq(seq)
fmt.Println(seq)
```

## Your task

Work out what it prints before you look at the answers.
