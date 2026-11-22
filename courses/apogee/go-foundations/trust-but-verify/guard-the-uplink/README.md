---
slug: guard-the-uplink
title: Guard the Uplink
type: code_tests
position: 8
difficulty: 7
topics: testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-13T11:32:35Z
---
Mirela's sign-off for the uplink: your checks for `verify` must catch all six bugs she has planted in copies of it.

A forgery should differ from a good command in exactly one way, so that when a check fails, it can fail for only one reason:

```go
good := Command{
	SatID: "furrow-a", Seq: 12, Op: "PING",
}
replay := good  // a copy...
replay.Seq = 11 // ...with exactly one change
```

## Your task

Write `checkVerify` in `checks.go`, with one subtest per case and `t.Errorf` for every wrong answer:

1. `labCommand`, signed with `labKey` by `sign`: valid.
2. Four forgeries of it, each changing one thing: the command, the signature's last hex digit, the key handed to `verify`, or an empty signature.
