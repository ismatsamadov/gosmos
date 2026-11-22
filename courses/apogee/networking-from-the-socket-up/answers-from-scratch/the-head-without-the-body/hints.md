`req.Method` says whether the request is `HEAD`. The answer `recording` builds already has the right status and fields for GET; for HEAD only its `body` changes.

Drop the body just before the function returns, after `a` has its final value, so the 404 for a pass that isn't kept loses its body too.

If Go's client sees a `ContentLength` of 0 for HEAD, the data was dropped before `Content-Length` was computed from it. If the raw connection still gets 4096 bytes, the body is never set to nil.
