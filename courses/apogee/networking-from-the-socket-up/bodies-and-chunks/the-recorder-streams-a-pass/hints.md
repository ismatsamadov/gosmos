`strings.CutPrefix(req.Target, "/recordings/")` gives the pass ID and whether the target started that way; both must hold before anything else is looked at.

`checkCRC(req.Body, req.Header, req.Trailer)` returns the CRC to answer with, or a `*RequestError` that `failure` turns into the 400. `s.Store.Put` returns an error whose text is already `<pass ID> already stored`.

Build the 201 with `textAnswer` and the `fmt.Sprintf` message, then append `Field{"Location", req.Target}` and the CRC as `%08x` to its fields. If the wrong-CRC pass appears in `/recordings`, the body is put in the store before its CRC is checked.
