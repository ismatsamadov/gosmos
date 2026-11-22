---
slug: live-in-the-fishbowl
title: Live in the Fishbowl
type: code_tests
position: 8
difficulty: 9
topics: errors, networking, strings, time
created: 2026-11-13T06:15:46Z
modified: 2026-11-13T16:03:24Z
---
Friday 16 June, 07:52: Furrow-A rises over Ystmark, and Ingrid watches its pass with the whole live view: every frame in order and once, and when the frames stop, a question to the recorder.

`Watch` brings the week's pieces together. The given `receive` waits up to `l.Quiet` for the recorder's next datagram, skipping strays and any over 512 bytes; when it fails with a deadline, the pass has gone quiet. `seqOf` reads a line's sequence number, and a datagram holds a frame a line:

```go
for _, line := range strings.Split("a\nb", "\n") {
	fmt.Println(line)
}
// a, then b
```

## Your task

`Watch` in `watch.go` prints whole datagrams until anything goes wrong. Make it:

1. Call `receive` until it fails with `os.ErrDeadlineExceeded`, returning any other error.
2. Print each line of each datagram whose `seqOf` number is accepted by `t`, one `Tracker` for the whole pass.
3. Then add one to `l.asked`, `ask` the recorder that question with `l.Wait` and `l.Tries`, and print `lastLine(t, answer)`, the answer being `silent` when every try timed out.
