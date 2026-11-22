---
slug: quick-check-shutdown
title: Quick Check: Shutdown
type: quiz
position: 3
difficulty: 9
topics: context, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Wednesday's rehearsal starts with Ingrid's go/no-go, and today its questions are about stopping. This was the first draft of `main`, written before `run`:

```go
go func() {
	<-ctx.Done()
	srv.Shutdown(context.Background())
}()
err := srv.Serve(ln)
log.Println("bookingd stopped:", err)
// main returns here, and the process ends
```

## Your task

Answer three questions: what `Shutdown` does with a request still being answered, why this draft can end the process too soon, and what `signal.NotifyContext` gives `main`.
