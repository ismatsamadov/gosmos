`netip.ParseAddrPort` reads `10.20.0.31:80` and `[::1]:443` alike, and its `Addr()` is the IP. A parse error is a refusal too: a check that can't tell must say no.

A `netip.Addr` answers each question itself, for IPv4 and IPv6 both: `IsLoopback`, `IsPrivate`, `IsLinkLocalUnicast` and `IsUnspecified`.

Wrap `errInside` with `%w` and add the address after it, so `errors.Is` still finds it under the dialer's own error, all the way up to the client's.
