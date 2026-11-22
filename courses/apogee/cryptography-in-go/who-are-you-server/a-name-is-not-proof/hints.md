`cert.Verify` takes one `x509.VerifyOptions` value and returns the chains it found and an error. You need only the error, so assign the chains to `_`.

Fill five fields of the options: `Roots`, `Intermediates`, `DNSName`, `KeyUsages` as a one-element `[]x509.ExtKeyUsage`, and `CurrentTime`. The name check on `CommonName` can go: `DNSName` checks the name properly.

If Mirela's certificate still passes, `trusted` still returns nil without asking `Verify`: return `Verify`'s error and nothing else.
