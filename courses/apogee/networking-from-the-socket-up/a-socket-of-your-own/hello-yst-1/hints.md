`net.DialTimeout` takes the same network and address as `net.Dial`, then a `time.Duration`: `time.Second` is one second.

After the error check, defer the connection's `Close`, then pass the connection to `handshake`, which returns the name, a `Welcome` and an error.

Build the line with `fmt.Sprintf` from the name and the Welcome's `FwMajor`, `FwMinor` and `LastSeq`. If the YST-2 case comes back as YST-1, the name was typed in rather than taken from `handshake`.
