Check the issuer first, as `newRoot` checks its key, with `errors.New` for the error. The template is an `x509.Certificate` with `SerialNumber` (`big.NewInt(serial)`), `Subject` (a `pkix.Name` with the host as `CommonName` and `Organization`), `DNSNames`, `NotBefore`, `NotAfter`, `KeyUsage` and `ExtKeyUsage`. Ninety days later is `now.AddDate(0, 0, 90)`.

`x509.CreateCertificate` takes `rand.Reader`, the template, the parent, the new certificate's public key and the signing key, and returns the certificate's bytes and an error; `x509.ParseCertificate` turns the bytes into the certificate to return.

If the chain doesn't verify, check the parent and the signer: the CA's certificate and the CA's key, never the template itself. If client authentication still verifies, `ExtKeyUsage` holds more than `ExtKeyUsageServerAuth`.
