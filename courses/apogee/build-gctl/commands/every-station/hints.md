Check `len(args)` before the loop, so that an argument gets only its error; its text is exactly `stations takes no arguments`, and `Run` adds the `error:` in front.

A small helper that returns the end of a station's line keeps `cmdStations` short: an empty string for a station that isn't stowed, ` stowed` when its `StowReason` is empty, and ` stowed: ` followed by the reason otherwise. The next lesson can use it too.

Each line is the code, a space, the name, then whatever the helper returns, and a newline. The table entry's `run` is `cmdStations`.
