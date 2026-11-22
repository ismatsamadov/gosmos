---
slug: remember-the-status
title: Remember the Status
type: code_tests
position: 3
difficulty: 7
topics: http, interfaces, structs
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Yusra wants a log line per request, saying what each got. Nkem asks first: "A handler writes 404. Who else sees it?" Nobody: a `ResponseWriter` passes it on and remembers nothing.

Embed an interface in a struct, and the struct gets its methods, passed to the value inside. A method of the same name on the struct runs instead, free to watch the call before passing it on:

```go
type counter struct {
	io.Writer
	n int
}

func (c *counter) Write(p []byte) (int, error) {
	c.n += len(p)
	return c.Writer.Write(p)
}
```

## Your task

Make `statusRecorder`'s `WriteHeader` record the first status in `Status`, and its `Write` record `200` if no status came first and add the bytes written to `Bytes`; both still pass the call on. Then add `Unwrap() http.ResponseWriter`, returning the writer inside.

> [!NOTE]
> `http.NewResponseController` calls `Unwrap` to reach the real writer's `Flush`, which sends what's written so far at once. Without `Unwrap`, the wrapper hides it.
