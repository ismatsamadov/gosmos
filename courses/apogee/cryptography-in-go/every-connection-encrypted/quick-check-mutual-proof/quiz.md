## On that server, a client sends a certificate it made itself, named uplinkd. What name does the draft read?

- [x] `uplinkd`, because `PeerCertificates` holds whatever the client sent
  > Nothing checked the certificate, so its name is only the client's claim, and the draft believes it.
- [ ] None: the handshake refuses a certificate no root signed
  > `RequestClientCert` asks but never verifies, so the handshake succeeds and the certificate arrives unchecked.
- [ ] None: `PeerCertificates` stays empty until a certificate is verified
  > `PeerCertificates` holds what was sent, verified or not; it is `VerifiedChains` that stays empty.

## uplinkd's verified certificate reaches a route that clears a station's whole queue. What decides whether it may?

- [ ] The handshake, since a verified client may do anything
  > The handshake proves who is calling; it says nothing about what that caller may do.
- [x] The handler, by checking the caller's name against what the route allows
  > A verified certificate is an identity, nothing more; each route still decides, as Ystmark's lets only uplinkd queue.
- [ ] `ClientCAs`, which lists what each client may do
  > `ClientCAs` holds the roots that client certificates must chain to, not any permissions.
