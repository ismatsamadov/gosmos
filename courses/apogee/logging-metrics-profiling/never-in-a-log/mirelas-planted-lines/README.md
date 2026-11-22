---
slug: mirelas-planted-lines
title: Mirela's Planted Lines
type: code_tests
position: 7
difficulty: 8
topics: http, logging, security
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Friday 1 December, 16:00. Before she signs off, Mirela plants twelve lines with a secret everywhere she found one this week: a header dump, a key inside an error, a `Key`, a password in a group, a token in the message, Hollis's address under `email`, and more. No secret may reach a disk, and all twelve lines still must.

One `ReplaceAttr` can apply every rule of the week, since the handler calls it for each attribute, the nested ones and the message included. `Key` and `APIKey` have logged themselves by then, and a type switch on the value finds the rest:

```go
switch v := a.Value.Any().(type) {
case time.Duration:
	return msAttr(a.Key, v)
}
return a
```

## Your task

Give `safeLogger`'s JSON handler a `ReplaceAttr` that runs `redact` first, then hands an `http.Header` to `safeHeader`, a string under the key `email` to `maskEmail`, and every other string and every error's text to `scrub`.

> [!NOTE]
> Values of any other type, such as the time, the level and a pass's numbers, must come back as they went in.
