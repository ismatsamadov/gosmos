---
slug: close-every-body
title: Close Every Body
type: code_tests
position: 4
difficulty: 6
topics: defer, http
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
After a week of polls, the poll's machine held hundreds of unfinished connections, and Tomás's weather service stopped taking new ones. Every answer whose body nobody closed was still holding a connection open.

A response's body keeps its connection busy until it is closed. Deferring the close right after the error check of `Do` closes it on every return that follows: an early one, a failed decode or a success. Files follow the same pattern:

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close()
```

## Your task

`Reading` closes the body only after a successful decode, so a 404, a 503 or broken JSON leaves it open. Defer `resp.Body.Close()` right after the error check of `Do`, and remove the call near the end.

> [!TIP]
> Close a body even when nothing reads it: until it is closed, its connection stays busy.
