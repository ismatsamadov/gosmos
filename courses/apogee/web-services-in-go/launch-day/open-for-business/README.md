---
slug: open-for-business
title: Open for Business
type: code_tests
position: 8
difficulty: 10
topics: context, errors, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T16:20:02Z
---
Tuesday 1 February, 14:05 UTC. Furrow-A is over Salar Quieto, customers are making bookings, and the guide's sample request says `Authorization: Token` where it should say `Bearer`. The guide is built into the service, so the service restarts now.

A restart that drops nothing stops in order: drain, shut down, then stop the notifier, which saves every delivery still in progress, attempts and all, to the outbox. The new process starts its notifier from that outbox before it serves anyone. `errors.Join` keeps every error that stopping returns and leaves out the nils:

```go
late := errors.New("shutdown: timed out")
err := errors.Join(late, nil)
fmt.Println(err) // shutdown: timed out
```

## Your task

`run` never starts or stops the notifier. Call `app.Notifier.Start()` before serving, and `app.Notifier.Stop()` once the server has shut down or failed, on every return, joining its error to the one `run` returns with `errors.Join`.
