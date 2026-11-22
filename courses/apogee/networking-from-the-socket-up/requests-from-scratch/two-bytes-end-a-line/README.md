---
slug: two-bytes-end-a-line
title: Two Bytes End a Line
type: code_tests
position: 3
difficulty: 8
topics: errors, http, streams
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T16:42:55Z
---
Mirela, Apogee's security engineer, reviews the station server. Her first rule: a line that doesn't end in CR LF, holds a CR anywhere else, or never ends is refused, never guessed at.

HTTP lets a server accept an LF alone; Mirela doesn't. The server reads through `bufio.NewReaderSize(conn, 1024)`, and `ReadSlice('\n')` returns a line inside that buffer: one over 1024 bytes, CR LF included, gives `bufio.ErrBufferFull`. `ReadString` grows to hold anything, even 64 KiB with no end:

```go
br := bufio.NewReaderSize(strings.NewReader(
	strings.Repeat("a", 2000)), 1024)
line, err := br.ReadSlice('\n')
fmt.Println(len(line), err)
// 1024 bufio: buffer full
```

## Your task

Make `readLine` use `ReadSlice('\n')`: `errLineTooLong` for `bufio.ErrBufferFull`; a `RequestError` 400 `bad line ending` unless the line ends in `\r\n` with no other `\r`; `io.ErrUnexpectedEOF` for a stream that ends mid-line (an empty one stays `io.EOF`); otherwise the line without CR LF.

> [!WARNING]
> `ReadSlice`'s bytes belong to the reader and change at its next read: return a string, which copies them.
