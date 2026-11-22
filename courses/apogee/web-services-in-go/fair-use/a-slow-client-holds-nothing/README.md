---
slug: a-slow-client-holds-nothing
title: A Slow Client Holds Nothing
type: code_tests
position: 5
difficulty: 8
topics: http, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Brinewatch's client also sent its headers a byte at a time, and each of its connections sat open for minutes. Each open connection takes some of the server's memory and one of its limited connection slots, so enough slow clients could keep everyone else out.

An `http.Server` waits as long as a client likes unless you set limits: `ReadHeaderTimeout` for the headers, `ReadTimeout` for the whole request, `WriteTimeout` for the answer, and `IdleTimeout` for a kept-alive connection between requests. `MaxHeaderBytes` caps the headers' size. A slow client looks like this:

```go
conn, err := net.Dial("tcp", addr)
if err != nil {
	return err
}
conn.Write([]byte("GET / HTTP/1.1\r\n"))
// ...and then nothing, for minutes
```

## Your task

`newServer` sets only `Handler`. Also set `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout` and `IdleTimeout` from `t`'s `ReadHeader`, `Read`, `Write` and `Idle`, and `MaxHeaderBytes` to 16 KiB, 16384.

> [!NOTE]
> Go reads 4 KiB past `MaxHeaderBytes` before it answers 431, so the test sends a 32 KiB header.
