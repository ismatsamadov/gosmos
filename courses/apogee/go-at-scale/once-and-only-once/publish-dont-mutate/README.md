---
slug: publish-dont-mutate
title: Publish, Don't Mutate
type: code_tests
position: 6
difficulty: 9
topics: goroutines, mutexes, pointers
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Thursday 22 March. After Wednesday's mix-up Nkem reads the publisher, and finds it writing each new plan's contacts into the very schedule the test Board is reading.

To share a value that changes, publish new copies: build the new value whole, then `Store` a pointer to it in an `atomic.Pointer`. A reader that loaded the old pointer keeps a schedule nobody changes under it, and its next `Load` gets the new one, complete. The one rule: a value is never changed once published.

```go
var latest atomic.Pointer[Weather]
w := &Weather{Station: "CWR", WindKmh: 74}
latest.Store(w) // w never changes again
cur := latest.Load()
fmt.Println(cur.WindKmh) // 74
```

## Your task

`Update` raises the version of the current schedule and copies the contacts into it. Build a new `Schedule` instead, its `Version` one more than the current one's (1 when there is none) and its `Contacts` a slice of its own, from `slices.Clone`; then `Store` and return it.

> [!WARNING]
> Reusing the old slice, as `append(s.Contacts[:0], ...)` does, saves memory and changes the contacts under every reader of the old schedule.
