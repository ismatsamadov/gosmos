---
slug: ask-again
title: Ask Again
type: code_tests
position: 6
difficulty: 9
topics: errors, loops, networking, time
created: 2026-11-13T06:15:46Z
modified: 2026-11-13T06:15:46Z
---
Thursday, 11:20, Furrow-A's pass. Ingrid asked the recorder how it was going, the question was lost on the island's line, and after a second's wait the live view gave up. Ingrid: "Lost question. Ask again."

Over UDP, silence is the only sign of a loss, and it can't say whether the question or the answer went missing. So give each try a deadline, and when it passes, send the same question again: under the same number, so a late answer to the first try still counts. After a few tries, stop.

```go
pc.SetReadDeadline(time.Now().Add(time.Second))
_, _, err := pc.ReadFrom(buf) // nobody sends
late := errors.Is(err, os.ErrDeadlineExceeded)
fmt.Println(late) // true
```

## Your task

`ask` in `ask.go` sends question `n` and waits `wait` for its answer, once. Make it try up to `tries` times: when a try's wait ends in `os.ErrDeadlineExceeded`, send the question again and wait again. Return the answer, any other error at once, or the last try's error.
