`r.c` is the client, so `r.c.conn` is its connection and `r.c.Timeout` its limit. Set the deadline first, then return what the connection's own `Read` returns.

Every byte `c.br` takes from the connection now goes through `idleReader.Read`, the banner included, so `recv` needs no deadline of its own: remove its `SetReadDeadline` line.

If the slow STATUS still fails after 5 s, a deadline is still set once for the whole frame. If the silent simulator fails after 0 s, the deadline is set after the `Read` instead of before it.
