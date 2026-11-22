---
slug: wrap-it-keep-it
title: Wrap It, Keep It
type: code_tests
position: 5
difficulty: 5
topics: errors, formatting
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
The log line `duplicate frame` doesn't say whose frame. Adding the satellite and the sequence number takes a format string; keeping the duplicate recognizable for the pager takes the right verb.

With `%v`, `fmt.Errorf` copies only an error's text. With `%w`, it *wraps* the error: the message reads the same, but the original stays inside, where `errors.Is` finds it. `==` compares only the outer error:

```go
err := fmt.Errorf("CWR: %w", ErrStowed)
fmt.Println(err) // CWR: dish stowed

fmt.Println(errors.Is(err, ErrStowed)) // true
fmt.Println(err == ErrStowed)          // false
```

## Your task

In `ingest`, change one verb so the error still reads `thistle seq 4182: duplicate frame` and `errors.Is` finds `ErrDuplicate` in it.
