---
slug: close-the-session
title: Close the Session, Always
type: code
position: 8
difficulty: 5
topics: defer, errors
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Last night a pass at Cape Wren failed on wind, returned early and never closed its session. The dish stayed booked, and the next pass found it taken.

`defer` schedules a call to run when the surrounding function returns, whichever `return` it takes. Put it right after you acquire something, so the cleanup sits next to the setup:

```go
func countdown() {
	defer fmt.Println("AOS")
	fmt.Println("T-2")
	fmt.Println("T-1")
}
// prints T-2, T-1, AOS
```

## Your task

Make `runPass` close its session on every path: defer the close right after `openSession`, and remove the call at the end.
