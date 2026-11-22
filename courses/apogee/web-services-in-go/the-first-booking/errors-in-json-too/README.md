---
slug: errors-in-json-too
title: Errors in JSON Too
type: code_tests
position: 1
difficulty: 6
topics: errors, http, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Furrowlight's program reads JSON, not sentences, and half the service's errors are still plain text, such as `no station "ZZZ"`. Nkem has sent every handler's errors through one function, `writeError`, for you to finish.

An `*APIError` carries what the caller should see: a status, a code and a message. `errors.AsType` finds one in an error even when other errors wrap it with `%w`. Any other error is the service's own failure, and its text may say things no customer should read:

```go
_, err := os.Open("guide/missing.html")
err = fmt.Errorf("guide: %w", err)
pe, ok := errors.AsType[*fs.PathError](err)
fmt.Println(ok, pe.Path)
// true guide/missing.html
```

## Your task

`writeError` answers 500 to every error. Answer an `*APIError`, wrapped or not, with its status and an `ErrorBody` of its code and message, and any other error with 500 and:

```text
{"error":"internal","message":"internal error"}
```

> [!WARNING]
> A type assertion, `err.(*APIError)`, looks only at the outermost error and misses a wrapped one.
