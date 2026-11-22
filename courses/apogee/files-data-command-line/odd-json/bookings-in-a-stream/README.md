---
slug: bookings-in-a-stream
title: Bookings in a Stream
type: code_tests
position: 0
difficulty: 5
topics: errors, loops, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:19:20Z
---
At each signal stop, the van's recorder adds the newest bookings to a log, one JSON object per line. Paz wants the whole log read tonight, but `json.Unmarshal` refuses the file at its second line.

A `json.Decoder` reads a stream one value at a time, skipping the white space between values. At a clean end, `Decode` returns `io.EOF`; a value cut off partway gives `io.ErrUnexpectedEOF`.

```go
dec := json.NewDecoder(strings.NewReader("7\n40"))
for {
	var n int
	if err := dec.Decode(&n); err != nil {
		break // io.EOF, after 40
	}
	fmt.Println(n)
}
```

## Your task

Implement `decodeStream`: decode bookings from `r` until the stream ends. At `io.EOF`, return them with a nil error; on any other error, return the bookings decoded so far and that error.
