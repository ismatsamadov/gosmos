## Which name does this client check Lantern Atoll's certificate against?

- [ ] `10.4.0.7`, the address it dials
  > The address says where to connect; with `ServerName` set, the certificate must name `lan.apogee.internal` instead.
- [x] `lan.apogee.internal`, from `ServerName`
  > `ServerName` is the name the certificate must carry, whatever address the client dials.
- [ ] None, since a TLS 1.3 server is trusted already
  > The version says how the connection is encrypted, not who is at the other end.

## Without the `RootCAs` line, which roots would the client trust?

- [ ] None, so it would refuse every certificate
  > A nil `RootCAs` doesn't mean no roots: Go falls back to the system's own.
- [ ] Whichever root the server sends with its chain
  > A root the server sends proves nothing, since anyone can make one, as the impostor did.
- [x] The system's own roots, which don't include Apogee's
  > With `RootCAs` nil, Go trusts the roots of the computer it runs on, and none of them is Apogee's, so Lantern Atoll would be refused.

## What does `InsecureSkipVerify: true` turn off?

- [x] Every check of the server's certificate: chain, name and dates
  > The bytes are still encrypted, but to whoever answers, impostors included.
- [ ] The encryption, so the bytes travel in plain text
  > The connection stays encrypted; what is lost is knowing who is at the other end.
- [ ] The name check alone, keeping the chain and the dates
  > It skips the whole verification, not one part of it.
