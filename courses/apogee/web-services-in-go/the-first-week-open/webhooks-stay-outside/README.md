---
slug: webhooks-stay-outside
title: Webhooks Stay Outside
type: code_tests
position: 5
difficulty: 10
topics: errors, http, security
created: 2026-11-13T06:35:19Z
modified: 2026-11-13T16:15:49Z
---
Thursday, 15:00. Tomás, on call this week, finds Brinewatch's events in his weather service's log. Brinewatch gave Yusra its receiver's address as it reads inside Brinewatch's own network, `http://10.20.0.31/apogee-events`, and inside Apogee's network 10.20.0.31 is the weather service. Mirela: "Next time, someone types an address inside our network on purpose."

A webhook may only go out to the internet: never to a loopback address (`127.x`, `::1`), a private one (`10.x`, `172.16.x` to `172.31.x`, `192.168.x`, `fc00::/7`), a link-local one (`169.254.x`, `fe80::/10`) or an unspecified one (`0.0.0.0`, `::`, which reach this very machine). Any name can resolve to them, so `webhookClient`'s dialer checks in its `Control` hook, which gets the very address it is about to connect to.

```go
ip, err := netip.ParseAddr("172.20.4.1")
fmt.Println(ip.IsPrivate(), err) // true <nil>
```

## Your task

`checkDial` allows every address. Parse `address`, an IP and a port, and return an error wrapping `errInside` for any of the addresses above, and for one that doesn't parse.
