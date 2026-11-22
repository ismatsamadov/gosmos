---
slug: string-is-not-enough
title: String Is Not Enough
type: code_tests
position: 2
difficulty: 7
topics: interfaces, json, logging, security
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Before Rotation Day you gave `Key` a `String` method that hides its secret from fmt: it prints `k-2028-10 (secret hidden)`. Mirela logs the lab bench's copy of `k-2028-10`, the master key in use since then, through the JSON handler, and reads its secret in base64.

The JSON handler never calls `String`. It writes a struct through its `LogValue` or, without one, as `encoding/json` would, every exported field included. The text handler formats a struct with fmt, which does call `String`. A `Shift` whose `String` returns its name:

```go
t.Info("on call", "shift", s) // text handler
// level=INFO msg="on call" shift="Mirela Dunn"
j.Info("on call", "shift", s) // JSON handler
// {"level":"INFO","msg":"on call","shift":
//  {"From":"2028-11-20T08:00:00Z",
//  "Name":"Mirela Dunn"}}
```

## Your task

Give `Key` a `LogValue` method that returns what its `String` method prints, as a string value.

> [!TIP]
> Return `slog.StringValue(...)`: a `slog.AnyValue(k)` would send slog back to `k`'s own `LogValue`.
