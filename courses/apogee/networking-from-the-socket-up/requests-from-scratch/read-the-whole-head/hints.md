Start with `readLine`: its `errLineTooLong` here means the request line was too long, so return the 414. Then `parseRequestLine` and `readHeader`, returning their errors unchanged; the server needs a deadline's error as it is to answer 408.

Count the Host field's values with `len(h["host"])`: 0 and more than 1 are the two 400s. Then check the method against GET, HEAD and PUT exactly, letter case included.

Last, `contentLength(method, h, max)` gives the length or the refusal. If two Host lines get through, the check used `h.Get`, which joins both values into one non-empty string.
