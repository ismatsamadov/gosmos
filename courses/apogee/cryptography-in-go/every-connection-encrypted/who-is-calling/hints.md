Look before you index: `r.TLS` may be nil, and `VerifiedChains` may be empty. Only when both hold is `r.TLS.VerifiedChains[0][0]` the client's certificate, and its `Subject.CommonName` its name.

In the `POST /uplink` handler, call `caller(r)` before reading the body, and refuse unless it returns true and `uplinkd`. `return` right after `http.Error`, so nothing reaches the queue.

If passplan still queues its command, the check comes after `q.Add` or misses its `return`. If `TestCaller` fails for the certificate that was never verified, `caller` reads `PeerCertificates`.
