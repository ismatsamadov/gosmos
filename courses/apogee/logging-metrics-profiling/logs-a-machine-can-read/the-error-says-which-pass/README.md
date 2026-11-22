---
slug: the-error-says-which-pass
title: The Error Says Which Pass
type: code_tests
position: 4
difficulty: 6
topics: errors, logging
created: 2026-11-07T07:51:36Z
modified: 2026-11-07T07:51:36Z
---
Yusra, on Thursday: "TL;DR: Cirrine still asks why its 01:55 command never left the ground." At 02:02:39 uplinkd knew which pass had failed: its error read `pass cirrine-1/LAN/16402: context deadline exceeded`, and the old line threw it away.

An error is a value like any other: `slog.Any("err", err)`, or `"err", err` among the pairs, writes its whole text, every wrapping included. Log it only when there is one, since a nil error writes `"err":null`:

```go
err := fmt.Errorf("booking bk-5499: %w",
	errors.New("not found"))
l.Error("request failed", "err", err)
// {"level":"ERROR","msg":"request failed",
//  "err":"booking bk-5499: not found"}
```

## Your task

`logResult` writes INFO `pass flown` with `pass_id` whatever happened. Keep that for a nil error; otherwise write ERROR `pass failed` with `pass_id` and `err`.

> [!WARNING]
> A flown pass with `"err":null` still has an `err` key, and a search for lines that have one would count it as a failure.
