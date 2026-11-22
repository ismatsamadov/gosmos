Start from what the chapter built: a client config with `RootCAs`, `ServerName`, `Certificates`, `MinVersion` and `Time`, inside an `http.Client` with a `Timeout`; a server config with `Certificates`, `MinVersion`, `ClientAuth`, `ClientCAs` and `Time`.

The server's `tls.Certificate` holds `leaf.Raw` then `inter.Raw`, with `key`. Both configs set `MinVersion` to `tls.VersionTLS13`, and the server's `ClientAuth` is `tls.RequireAndVerifyClientCert`.

Read which case fails: an impostor station getting through means the client's roots or name are unset; an impostor uplinkd getting through means the server's client checks are; the TLS 1.2 client or station getting through means the server's or the client's `MinVersion` is missing.
