---
slug: errors-that-carry-the-status
title: Errors That Carry the Status
type: code_tests
position: 7
difficulty: 6
topics: errors, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Ingrid wants the poll to tell a silent weather mast, the pole of instruments that measures a station's wind, from a typo in a station's code. A 503 means wait for the next reading and a 404 is a bug to fix, yet both reach the poll as one line of text.

An error can be a value of your own type that carries data: `StatusError`, in `weather.go`, holds the status and the start of the body. Wrapped with `%w`, it travels inside the message, and `errors.AsType` gets it back:

```go
_, err := os.ReadFile("missing.txt")
pe, ok := errors.AsType[*fs.PathError](err)
fmt.Println(ok, pe.Op, pe.Path)
// true open missing.txt
```

## Your task

`Reading` turns any status but 200 into plain text. Instead, read at most 512 bytes of the body through `io.LimitReader`, and wrap a `*StatusError` holding the code, the status and that text, trimmed of white space, with `%w`, so the typo's error reads:

```text
weather for CRW: 404 Not Found: no station "CRW"
```
