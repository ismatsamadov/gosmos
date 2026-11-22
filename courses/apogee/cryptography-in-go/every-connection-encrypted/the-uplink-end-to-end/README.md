---
slug: the-uplink-end-to-end
title: The Uplink, End to End
type: code_tests
position: 8
difficulty: 10
topics: crypto, http, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T19:06:42Z
---
Friday 6 October. Every station now speaks TLS 1.3 with both ends proven, and today Cirrine-1's sealed commands cross it to Ystmark. The station queues them unread: it holds no key that opens them, and only the satellite can.

Two configs carry everything this chapter required. The sealed line hides its command from everyone on the way, the station included; TLS adds a proven station and a proven uplinkd, and hides even the line's header from the networks in between:

```go
c := uplinkClient(roots, "yst.apogee.internal",
	cert, now)
err := send(ctx, c, ystURL, sealedLine)
fmt.Println(err) // <nil>: Ystmark queued it
```

## Your task

1. `uplinkClient`: trust only `roots`, check `host`, send `cert`, TLS 1.3 only, `Time: now` and a 2-second timeout.
2. `stationServer`: present `leaf` then `inter` with `key`, TLS 1.3 only, require clients that chain to `roots`, and `Time: now`.
