`crypto/tls` has its own dialer: a `tls.Dialer` holds a `Config`, and its `DialContext` takes the same context, network and address as the plain dialer's.

Make a `tls.Dialer` whose `Config` is `cfg`, and return what its `DialContext` returns for `"tcp"` and `addr`: a connection whose handshake is already done, or the error that stopped it.

If the silent listener's test says it got a connection, the handshake hasn't happened yet, as with a bare `tls.Client`; if it says the dial was still waiting, the handshake ignores `ctx`, as `tls.Dial` does. If the echo comes back empty, the line went out unencrypted, and the controller hung up on it.
