`pem.Decode(data)` returns the first block and the rest; check that the block isn't nil before passing `block.Bytes` to `x509.ParseCertificate`.

The fields are `cert.Subject.CommonName`, `cert.Issuer.CommonName`, `cert.SerialNumber`, `cert.NotBefore` and `cert.NotAfter` (each formatted with `time.DateOnly`), `cert.DNSNames` and `cert.IsCA`; `fmt.Println` prints the slice in brackets.

If a line still differs from the expected output, compare its label and spacing: one space after each colon, and the word `to` between the two dates.
