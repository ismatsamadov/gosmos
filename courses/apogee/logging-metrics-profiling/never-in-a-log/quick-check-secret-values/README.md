---
slug: quick-check-secret-values
title: Quick Check: Secret Values
type: quiz
position: 3
difficulty: 7
topics: logging, security
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Monday, 16:00: Yusra and Hollis have replaced Furrowlight's key (Hollis: "NEW KEY. WHY."). Before more code changes, Mirela writes five lines on the whiteboard.

Here `k` is a `Key` with its `String` method and no `LogValue`, and the text of `err` holds Furrowlight's old API key:

```go
// blank: the attribute "key" gets "[redacted]"
o := &slog.HandlerOptions{ReplaceAttr: blank}
j := slog.New(slog.NewJSONHandler(os.Stdout, o))
j.Info("command sealed", "seq", 4688, "key", k)
j.Error("request failed", "err", err)
```

## Your task

Answer Mirela's three questions: how the JSON handler writes `k` without `blank`, which attributes `blank` is called for, and what it does about the API key inside `err`.
