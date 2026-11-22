---
slug: lengths-you-can-trust
title: Lengths You Can Trust
type: code_tests
position: 5
difficulty: 7
topics: errors, networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
During Sunniva's cable test on Thursday afternoon, a damaged length read `ff ff`, and dishd sat waiting for 65,535 bytes that would never come. A frame has 9 to 1,024 bytes after its length, so any other length means the stream is out of step.

Check a length before you trust it: before reading one more byte, and before making a slice that big. Wrap the error with the value, so the log says which length it was and `errors.Is` still finds the error:

```go
ErrTooCold := errors.New("too cold")
err := fmt.Errorf("%w %d", ErrTooCold, -41)
fmt.Println(err, errors.Is(err, ErrTooCold))
// too cold -41 true
```

## Your task

`readPayload` reads any length it is given. Replace the `TODO` so a length under 9 or over 1024 returns `ErrBadLength` wrapped with the length, as `bad length 65535`, without reading further.
