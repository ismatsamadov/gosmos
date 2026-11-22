---
slug: closed-is-not-failed
title: Closed Is Not Failed
type: code_tests
position: 1
difficulty: 9
topics: errors, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Tuesday's rehearsal log says `bookingd failed: http: Server closed` every time someone stopped the service on purpose, and the pager, the alarm service that wakes whoever is on call, woke Nkem twice.

`srv.Serve(ln)` returns only when serving is over, and always with an error. After `Shutdown` or `Close` that error is `http.ErrServerClosed`: the server stopped because it was told to. Anything else is a real failure, such as a listener that broke:

```go
err := srv.Serve(ln) // until Shutdown
fmt.Println(err)
// http: Server closed
```

## Your task

`serve` returns whatever `Serve` returns. Return nil when the error is `http.ErrServerClosed`, found with `errors.Is`, and the error itself otherwise.

> [!TIP]
> Never compare an error's text: the text is for people, and it may change in the next release of Go.
