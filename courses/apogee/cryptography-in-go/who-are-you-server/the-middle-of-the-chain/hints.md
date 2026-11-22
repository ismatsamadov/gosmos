`x509.NewCertPool()` makes an empty pool, and its `AddCert` method adds one certificate; a `for` loop over `chain` adds them all.

Keep `Roots`, `DNSName` and `CurrentTime` as they are, and add two fields: `Intermediates` with your pool, and `KeyUsages` with a one-element `[]x509.ExtKeyUsage` holding `x509.ExtKeyUsageServerAuth`.

If Ystmark's certificate with its chain still fails, check that the pool you filled is the one in `Intermediates`: a pool made but never passed changes nothing.
