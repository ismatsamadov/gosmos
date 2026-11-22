---
slug: redact-by-name
title: Redact by Name
type: code_tests
position: 4
difficulty: 7
topics: logging, security, strings
created: 2026-11-10T15:36:38Z
modified: 2026-11-12T11:41:59Z
---
Mirela searches the old lines for more: attributes named `token`, `password` or `Authorization` turn up in records nobody planned, and she wants them blank by rule, whoever writes the next one.

A handler calls its `ReplaceAttr` for every attribute it writes that isn't a group, the built-in ones for the time, the level and the message too, with the names of the groups around it, and writes whatever comes back:

```go
show := func(gs []string, a slog.Attr) slog.Attr {
	fmt.Println(gs, a.Key)
	return a
}
h := slog.NewJSONHandler(io.Discard,
	&slog.HandlerOptions{ReplaceAttr: show})
p := slog.Group("pass", "station", "LAN")
slog.New(h).Info("pass started", p)
```

It prints the groups and the key of each:

```text
[] time
[] level
[] msg
[pass] station
```

## Your task

`redact` returns every attribute as it came. Make it replace the value of any attribute keyed `api_key`, `authorization`, `password`, `secret` or `token`, in any letter case and at any depth, with `[redacted]`, keeping the key.

> [!WARNING]
> Match whole keys: `token_count` is not on the list, and its number must stay.
