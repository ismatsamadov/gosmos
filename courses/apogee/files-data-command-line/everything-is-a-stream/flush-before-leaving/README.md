---
slug: flush-before-leaving
title: Flush Before Leaving
type: code_tests
position: 7
difficulty: 3
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-07T19:09:32Z
---
The report looked complete in the van's program, but its last lines never reached the stick. They were still waiting in a buffer when the helper returned.

`bufio.Writer` holds small writes and sends larger blocks to its destination. A write can succeed by entering that buffer; `Flush` sends the remaining bytes and reports any delayed failure:

```go
b := bufio.NewWriter(destination)
_, err := fmt.Fprintln(b, "ready")
if err != nil {
    return err
}
return b.Flush()
```

## Your task

Implement `bufferedReport`: make a buffered writer, write each supplied line with `fmt.Fprintln`, and flush before returning success. Return errors from either writing or flushing. Do not close the destination, which belongs to the caller.
