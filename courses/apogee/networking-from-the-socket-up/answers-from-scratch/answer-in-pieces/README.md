---
slug: answer-in-pieces
title: Answer in Pieces
type: code_tests
position: 5
difficulty: 10
topics: channels, http, streams
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T19:06:41Z
---
Monday evening. For Ingrid to hear of each recording as it lands, the collector will keep one request open on the station server's stream of arrivals, the feed at `/feed`.

An answer of unknown length can be chunked, like an upload: each piece is its size in lowercase hex, CR LF, its bytes and CR LF, and a chunk of size 0 and a blank line end it. Flush each chunk once written, or it waits in the buffer:

```go
fmt.Fprintf(os.Stdout, "%x\n", 4096) // 1000
```

## Your task

1. `writeChunk`: write `p` as one chunk and flush it; an empty `p` writes nothing.
2. `endChunks`: write the last chunk, `0\r\n\r\n`, and flush.
3. `(s *Server) feed`: send each line from `lines`, with a newline, as one chunk, returning the first error; when `lines` closes, call `endChunks`.

> [!NOTE]
> `route` now answers `GET /feed`: the given `send` writes its head, then hands your `feed` each new recording's line.
