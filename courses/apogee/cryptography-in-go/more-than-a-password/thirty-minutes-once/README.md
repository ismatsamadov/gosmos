---
slug: thirty-minutes-once
title: Thirty Minutes, Once
type: code_tests
position: 1
difficulty: 9
topics: crypto, maps, time
created: 2026-11-13T06:18:23Z
modified: 2026-11-13T16:11:48Z
---
Thursday 26 October. The oldest link in the reset file was emailed in March 2027 and still works, and a link Nkem used in May works a second time. Mirela: "A reset link is a password sent by email. It works once, for half an hour, then never."

`Start` already writes each reset's `Expires`, 30 minutes after it starts on the console's clock, `r.now`, but nothing reads it. A link is good from its start up to, not including, `Expires`. Once a token is accepted, its reset is deleted, so the link can't set a password twice:

```go
due := map[string]int{"pipit-1": 4182}
delete(due, "pipit-1")
_, ok := due["pipit-1"]
fmt.Println(ok) // false
```

## Your task

Make `Redeem` refuse a reset when `r.now()` is at or past its `Expires`, and delete the reset from `r.pending` once it accepts the token. Every refusal returns `errInvalidReset` itself and changes nothing, so a wrong guess can't spend a real link.
