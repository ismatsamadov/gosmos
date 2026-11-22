---
slug: is-it-a-timeout
title: Is It a Timeout?
type: code_tests
position: 1
difficulty: 7
topics: errors, networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Tuesday afternoon, Ingrid asks for a log she can read at once: `YST-1: no answer in 5s` when the line goes quiet, `YST-1: closed` when the controller hangs up, and exact words for everything else.

A read past its deadline returns an error matching `os.ErrDeadlineExceeded`, and `errors.Is` finds it however deeply it's wrapped. A stream that ends shows as `io.EOF`, or `io.ErrUnexpectedEOF` if `io.ReadFull` was partway through. A `time.Duration` prints itself as `50ms` or `5s`:

```go
err := fmt.Errorf("reading STATUS: %w", io.EOF)
fmt.Println(err == io.EOF)          // false
fmt.Println(errors.Is(err, io.EOF)) // true
fmt.Println(5 * time.Second)        // 5s
```

## Your task

Make `describe` return `YST-1: no answer in 50ms` for a deadline error, using `timeout`, and `YST-1: closed` for `io.EOF` or `io.ErrUnexpectedEOF`, wrapped or not. Other errors keep their own text.

> [!WARNING]
> Don't search the text for "timeout": `iotest.ErrTimeout`, a pretend error from Go's `testing/iotest`, says just that and is no deadline; `context.DeadlineExceeded` is a context running out.
