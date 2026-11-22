`w.Dial(ctx, c.Station)` returns a `net.Conn`; give it the pass's context, the one that ends at LOS. Put one `bufio.Reader` around the connection for the whole call: its `ReadString('\n')` returns the whole line however many pieces it came in, with the `\n` and any `\r` before it.

Trim the `\n`, then a `\r` if one is left, with `strings.TrimSuffix`, and compare what remains with the station's code followed by a space and `READY PROTO 1`. Keep in the pass whether the banner came right: its LOS line says `flown` only then.

Once the banner is right, wait for the pass's context to end, then `Close` the connection: the simulated world notes the instant you hang up, and the tests want each call's LOS. After a failed dial or a wrong banner there is nothing to hold, so hang up at once; the LOS line still waits for LOS.
