---
slug: encrypting-a-bare-socket
title: Encrypting a Bare Socket
type: code_tests
position: 7
difficulty: 9
topics: context, crypto, networking
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T19:06:42Z
---
Sunniva, about Ystmark's dish controllers, the computers that turn the dishes: "Snow 20 cm. Controllers still talk in the clear. Want them locked before the storm." They speak their own protocol of raw bytes, not HTTP, and only dishd, the dish service at Port Ellery, talks to them.

TLS wraps any connection. `tls.Client` puts it around a `net.Conn` you already have, but it handshakes only when first used, so a controller that never answers could hold its caller for as long as it likes. A `tls.Dialer` dials and handshakes in one call that obeys its context:

```go
tc := tls.Client(raw, cfg) // no handshake yet
st := tc.ConnectionState()
fmt.Println(st.HandshakeComplete) // false
```

## Your task

Make `dialController` dial with a `tls.Dialer` whose `Config` is `cfg`, so that when it returns, the handshake is done or has failed, within `ctx`. The test's stand-in controller echoes each line; another listener never answers at all.

> [!NOTE]
> `tls.Client` followed by `HandshakeContext(ctx)` works as well; then close the connection when the handshake fails.
