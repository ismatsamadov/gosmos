`WriteFile` takes a path, a byte slice and permission bits. Convert the supplied text instead of writing a fixed report.

The function already returns an error, so its result can be passed straight back to the caller.
