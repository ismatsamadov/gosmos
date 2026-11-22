---
slug: a-value-that-logs-itself
title: A Value That Logs Itself
type: code_tests
position: 1
difficulty: 7
topics: interfaces, logging, security
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Furrowlight gets a new key by 16:00, and Mirela wants that one safe wherever someone logs it: "Don't count on every log call remembering. Make the key refuse to be logged."

A type that implements `slog.LogValuer`, the one method `LogValue() slog.Value`, decides how it appears in every record: the JSON and text handlers call it before they write the value, inside groups and through `With` too. A `Station` can log as its code alone:

```go
func (s Station) LogValue() slog.Value {
	return slog.StringValue(s.Code)
}

l = l.With("site", lan) // lan is a Station
l.Info("pass started")
// {"level":"INFO","msg":"pass started",
//  "site":"LAN"}
```

## Your task

Give `APIKey` a `LogValue` method that returns its customer ID, a dot and `[redacted]`, as in `furrowlight.[redacted]`, or `[redacted]` alone when the key has no dot.

> [!NOTE]
> `LogValue` is for slog alone: `fmt.Println(k)` still prints the whole key.
