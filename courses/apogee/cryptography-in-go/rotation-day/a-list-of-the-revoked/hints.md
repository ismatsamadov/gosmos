`x509.CreateRevocationList` takes a source of randomness, `rand.Reader`, then a template, the CA's certificate and its key. The template is an `x509.RevocationList` with `Number` (a `*big.Int`), `ThisUpdate`, `NextUpdate` and `RevokedCertificateEntries`.

In `isRevoked`, take the steps in order: `x509.ParseRevocationList`, then the list's `CheckSignatureFrom(ca)`, then the dates, returning each failure as an error. Only then range over `RevokedCertificateEntries`.

Serial numbers are `*big.Int` values, compared with `Cmp`, never `==`. If the list still counts a second after `NextUpdate`, the date check is off by one: a list is current while `now` is neither before `ThisUpdate` nor after `NextUpdate`.
