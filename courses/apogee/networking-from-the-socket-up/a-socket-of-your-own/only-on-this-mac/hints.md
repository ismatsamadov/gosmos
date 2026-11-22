An address is a host, a colon and a port. Keep the host, `127.0.0.1`, and change only the port.

Port `0` doesn't mean port zero: it asks the system for whatever port is free right now, so two listeners never ask for the same one.

If the test says a second listener couldn't start, the port is still fixed. If it says the host is `::`, the host went missing: `:0` listens on every network.
