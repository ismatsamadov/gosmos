---
slug: a-datagram-from-the-hut
title: A Datagram from the Hut
type: code
position: 0
difficulty: 8
topics: networking, printing
created: 2026-11-13T06:15:45Z
modified: 2026-11-13T06:15:45Z
---
Monday 12 June, 06:30. During each pass the recorder will now send Port Ellery the satellite's readings of its own health, which ops calls telemetry, for the live view you're writing for Ingrid's screen. Pipit-2's 06:36 pass is the first test. Sunniva: "Recorder sending. Anyone listening?"

The readings travel as datagrams, over UDP: no connection, each datagram on its own. `net.ListenPacket("udp", "127.0.0.1:0")` opens a socket on a port the system picks; `WriteTo` sends one datagram to an address, and `ReadFrom` waits for one, copies it into your slice and returns its length and its sender:

```go
pc, _ := net.ListenPacket("udp", "127.0.0.1:0")
pc.WriteTo([]byte("ping"), pc.LocalAddr())
buf := make([]byte, 64)
n, from, _ := pc.ReadFrom(buf)
self := from.String() == pc.LocalAddr().String()
fmt.Println(n, self) // 4 true
```

## Your task

Replace the `TODO`: send `frame` from `rec` to `view.LocalAddr()`, read it on `view` into `buf`, setting `n`, and set `fromRecorder` to whether its sender is `rec.LocalAddr()`. It prints `70 bytes, from the recorder: true`, then the frame.

> [!NOTE]
> `WriteTo` returns once the datagram has left and never learns whether it arrived: sent where nobody listens, a datagram vanishes without an error.
