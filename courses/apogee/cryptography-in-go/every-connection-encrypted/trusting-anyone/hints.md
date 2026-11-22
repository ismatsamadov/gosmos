`RootCAs` is the pool of roots a client trusts, and `ServerName` the name the server's certificate must carry. Both are fields of the `tls.Config` inside the transport, where `InsecureSkipVerify` sits now.

Delete the `InsecureSkipVerify` line and set the two fields from `roots` and `host`, keeping `Time: now`. The timeout belongs to the `http.Client` itself, beside `Transport`.

If `TestYstmark` fails with `doesn't contain any IP SANs`, the name is missing; with `certificate signed by unknown authority` or `broken key size`, the roots are. If only `TestTimeout` fails, set `Timeout` to two seconds.
