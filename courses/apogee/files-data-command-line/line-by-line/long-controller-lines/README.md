---
slug: long-controller-lines
title: Long Controller Lines
type: code_tests
position: 2
difficulty: 4
topics: streams, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T16:34:04Z
---
On 2 July Paz moved La Grande slowly back and forth across a target while the controller logged where it pointed, many times a second: a calibration sweep. Its record, which Paz brought too, holds the whole sweep on one line of about 200 KiB, and the scanner stops there with an error.

A scanner's buffer must hold a whole line, newline included, and by default it grows to 64 KiB; a line that does not fit makes `Err` return `bufio.ErrTooLong`. `sc.Buffer(buf, max)` gives it a starting buffer and lets it grow to `max` bytes. Keep a bound: a damaged file without newlines should stop the scan, not fill memory.

```go
sc := bufio.NewScanner(os.Stdin)
buf := make([]byte, 0, 4096)
sc.Buffer(buf, 1024*1024) // up to 1 MiB
```

## Your task

At the TODO in `scanRecords`, let the buffer grow to `maxRecordBytes`, 256 KiB. The sweep must scan, and a longer record must still stop with `bufio.ErrTooLong`.

> [!WARNING]
> Call `Buffer` before the first `Scan`: once scanning has started, it panics.
