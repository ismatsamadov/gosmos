---
slug: trusting-anyone
title: Trusting Anyone
type: code_tests
position: 0
difficulty: 9
topics: crypto, http, security
created: 2026-11-07T07:46:46Z
modified: 2026-11-07T17:01:10Z
---
Monday 25 September. Tomás has shown that uplinkd's traffic to Ystmark and Lantern Atoll crosses another company's network as plain HTTP, and Mirela is reviewing its encrypted replacement. One line, added to get past an error, stops her: `InsecureSkipVerify: true`. "Then any server that answers is Ystmark."

The protocol under HTTPS, TLS, encrypts a connection and proves who is at the far end: the server presents its certificate, and the client checks that it chains to a root the client trusts and names the host it meant to reach. `InsecureSkipVerify` turns both checks off: any certificate, for any name, gets through. In production, the client reads its roots from `roots.pem`:

```go
data, err := os.ReadFile("roots.pem")
if err != nil {
	return err
}
roots := x509.NewCertPool()
roots.AppendCertsFromPEM(data)
```

## Your task

In `newClient`, replace `InsecureSkipVerify` with `RootCAs: roots` and `ServerName: host`: the tests dial 127.0.0.1, so the client needs the name to check. Give the `http.Client` a 2-second `Timeout`. The test plays an impostor with a certificate for example.com.

> [!NOTE]
> `Time: now` makes the client judge dates by the test's clock, Monday 25 September 2028.
