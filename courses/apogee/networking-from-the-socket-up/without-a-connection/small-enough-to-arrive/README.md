---
slug: small-enough-to-arrive
title: Small Enough to Arrive
type: code_tests
position: 2
difficulty: 8
topics: networking, slices, strings
created: 2026-11-13T06:15:46Z
modified: 2026-11-13T06:15:46Z
---
Tuesday 13 June, 08:06, Furrow-B's pass: the recorder fell nine seconds behind, then sent the nine frames it owed as one datagram of 656 bytes. The live view refused it, as it should, and Ingrid's screen skipped nine seconds.

A datagram bigger than one packet of the line, about 1,500 bytes, travels in pieces, and losing any piece loses all of it. So the recorder's limit is 512 bytes, well inside one packet, and every datagram stands alone: whole lines only, so one that goes missing costs its own frames and never half of another. Tomás: "A jigsaw in three envelopes: lose one and there's no picture."

```go
d := strings.Join([]string{"ab", "cde"}, "\n")
fmt.Println(len(d)) // 6: the newline counts
```

## Your task

The recorder packs its frames with `pack`, in `pack.go`, which puts every line in one datagram. Fill each datagram with whole lines in order, joined by `\n`, and start the next one whenever another line would take it past `maxDatagram`. No lines make no datagrams, and a single line over `maxDatagram` is an error wrapping `ErrTooBig`.
