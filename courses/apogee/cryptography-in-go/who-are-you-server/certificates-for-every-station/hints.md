`ed25519.GenerateKey(rand.Reader)` returns the public key, the private key and an error. Choose the function by `h.Client`: a variable can hold either `issueServer` or `issueClient`, since they take the same arguments.

`pem.EncodeToMemory` encodes one `pem.Block`: for a certificate, the type `CERTIFICATE` and the certificate's `Raw` bytes. Append the stations CA's block, made from `ca.Cert`, the same way. The key file is one `PRIVATE KEY` block holding what `x509.MarshalPKCS8PrivateKey` returns.

If `TestFailureLeavesNothing` finds files, `issueAll` saved some before every holder was issued: keep every file in memory first, and call `saveFile` only once the loop over `holders` has finished.
