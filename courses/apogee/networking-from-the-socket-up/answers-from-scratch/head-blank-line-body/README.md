---
slug: head-blank-line-body
title: Head, Blank Line, Body
type: code_tests
position: 1
difficulty: 9
topics: http, streams, time
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Monday morning. The collector, the program at Port Ellery that fetches each recording, will read these answers across the island's line, and Nkem wants each one sent exactly, in few writes.

After the status line come `Date`, the server's clock in UTC, the answer's fields, each a name, a colon, a space, its value and CR LF, then a blank line and the body. A `bufio.Writer` gathers it all and writes only when its 4,096 bytes fill or at `Flush`, which returns the first error any write met:

```go
bw := bufio.NewWriter(os.Stdout)
bw.WriteString("Wind 64\n")
bw.WriteString("Dish 1 stowed\n")
err := bw.Flush() // one write, both lines
```

## Your task

`writeResponse` writes only the status line. Through one `bufio.Writer`, write the status line, `Date` from `now` in UTC with the given `dateLayout`, each field in order, the blank line and the body, and return `Flush`'s error.

> [!WARNING]
> `dateLayout` prints `GMT` whatever the time zone: convert `now` with `UTC()` first.
