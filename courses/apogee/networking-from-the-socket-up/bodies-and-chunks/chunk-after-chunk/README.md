---
slug: chunk-after-chunk
title: Chunk After Chunk
type: code_tests
position: 3
difficulty: 9
topics: http, loops, streams
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Tuesday 6 June, 06:50. Sunniva: "Calm. Furrow-B at 07:15, first pass to stream in." The recorder streams it as it comes down, in chunks of whatever size it has ready.

A chunked body is a series of chunks: a size line, that many bytes, then CR LF. A chunk of size 0 ends it, and trailer fields follow, up to a blank line. Go's client sends 24 bytes like this:

```text
c\r\nfirst piece \r\n
c\r\nsecond piece\r\n
0\r\n
Apogee-Crc32: e2fc055d\r\n
\r\n
```

## Your task

`readChunked` returns nothing. Loop:

1. Read a size line with `readLine` and `chunkSize`; one too long for the reader is 400 `chunk line over 1024 bytes`.
2. Read the chunk and the two bytes after it with `io.ReadFull`; unless they are CR LF, refuse with 400 `bad chunk ending`.
3. At size 0, read the trailer with `readHeader` and return.

> [!WARNING]
> Never read a chunk with `readLine`: a recording holds CR and LF bytes.
