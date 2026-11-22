Start with the key: when its length isn't `ed25519.PrivateKeySize`, return an error made with `errors.New`, before anything calls `key.Public()`.

The template sets `SerialNumber` (`big.NewInt(1)`), `Subject` (a `pkix.Name` with `CommonName` and `Organization`), `NotBefore`, `NotAfter`, `IsCA`, `BasicConstraintsValid`, `MaxPathLen` and `KeyUsage`. `x509.CreateCertificate` takes `rand.Reader`, the template, the parent, the new certificate's public key and the signing key: for a root, the template is its own parent and the key signs its own public half.

If `NotAfter` comes out two days early, the ten years were counted in hours: `AddDate` counts calendar years, leap days included.
