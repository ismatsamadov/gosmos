`requireClients` sets two fields through its pointer, as `harden` set one: the policy, `ClientAuth`, and the pool that client certificates must chain to, `ClientCAs`.

The policy that both asks for a certificate and verifies it is `tls.RequireAndVerifyClientCert`. In `clientConfig`, `Certificates` is a slice of `tls.Certificate`: put `cert` in it.

If uplinkd's request fails with `certificate required`, `clientConfig` sends nothing yet. If a client that should be refused still gets through, check that `ClientAuth` verifies, not only requests.
