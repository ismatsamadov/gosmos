Keep, for each satellite, the masters it holds: at the window's start every one activated before it and still usable, then each one whose LOAD_KEY ran. Work out a pass's commands when it first takes the satellite's lock, so it never sends a LOAD_KEY that another pass has already had run.

The LOAD_KEY's plaintext joins with bars the operation, the master's ID, and the satellite's sign and seal keys from that master in lowercase hex; it is sealed like a Cirrine command, whatever the satellite. Mark the master held when its ACK comes, when a WELCOME's `last` shows it ran, or when a copy meets `replayed sequence number`.

`command failed authentication` or `unknown key` means a command was made with a master the satellite doesn't hold yet: Wednesday at Port Ellery starts with `k-2028-10` alone. If a LOAD_KEY ran twice, keep the one in flight at LOS beside its satellite, and settle it from the next WELCOME before working out that pass's commands.
