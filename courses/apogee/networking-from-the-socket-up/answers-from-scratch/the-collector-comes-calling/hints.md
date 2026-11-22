`strings.CutPrefix(req.Target, "/recordings/")` gives what follows the prefix and whether it was there; the target is a recording only when `validPassID` accepts what follows. Every other target is compared whole, so `/recordings?sat=furrow-a` is unknown.

A `switch` with no value tries its cases in order. Give each case a path and the methods it takes, then a case for the same path with any other method, which is the 405.

`refuse` drops the text for HEAD, so HEAD `/feed` gets its 405 and `Allow: GET` alone. If PUT `/recordings` comes back 404, a known path is answered as unknown: it takes GET and HEAD, so it's 405 with `Allow: GET, HEAD`.
