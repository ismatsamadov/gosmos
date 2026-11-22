`VerifyConnection` is a field of type `func(tls.ConnectionState) error`. Assign it a function literal that uses `current`, `ca` and `now` from `withRevocations`' parameters, so every handshake reads them afresh.

Inside, take the list from `current()`: nil means no list yet, so return nil. Otherwise call `isRevoked` with the list, `ca`, `cs.PeerCertificates[0]` and `now()`, and return its error when there is one.

If the second connection in `TestListChanges` still gets through, `current()` ran once, when `withRevocations` did, instead of inside the function; if only `TestResumedSession` fails, the check sits in `VerifyPeerCertificate`, which a resumed session skips. Build the refusal with `fmt.Errorf`, printing the serial number with `%v`.
