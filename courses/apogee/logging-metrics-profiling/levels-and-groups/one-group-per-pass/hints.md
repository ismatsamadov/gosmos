`slog.Group` takes the group's key first, `"pass"`, then the attributes inside it, each made with the function for its type: `slog.String` for the satellite and the station, `slog.Time` for the two times.

The four come from `p`'s fields `SatID`, `Station`, `AOS` and `LOS`, in that order; call `UTC()` on each time before you hand it to `slog.Time`.

If the test shows `pass.aos = "2028-11-07T18:55:09-07:00"`, a time went in with its own offset: give `slog.Time` the time's `UTC()`. If it shows `"pass":"cirrine-1/LAN/16402"`, the starter's string is still returned.
