A field is present when the header has its name as a key, whatever its value: `_, ok := h["transfer-encoding"]`. `h.Get` joins a repeated field's values with `", "`, so two lines of `chunked` read `chunked, chunked`, which ends in `chunked` but isn't `chunked` alone.

Check both fields first and refuse at once; then, with only `Transfer-Encoding`, its last coding, after the last comma and trimmed of spaces and tabs, then the whole value, then the method. With only `Content-Length`, or neither, `contentLength` already answers.

If the both-fields cases come back as a length of 46, `Content-Length` still wins when both are present, which is the trick itself. If `chunked, gzip` gets 501, only the whole value is checked. Both coding refusals quote the joined value with `%q`, and the method's name starts its 400, as in `HEAD has no body`.
