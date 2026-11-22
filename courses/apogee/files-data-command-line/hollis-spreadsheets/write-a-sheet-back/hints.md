Make one `csv.Writer` around `w` and call its `Write` method for each row, returning the error if one comes back.

After the loop, call `Flush`: a small sheet fits in the writer's buffer, and without it nothing reaches `w` at all.

`Flush` has no result, so finish with the writer's `Error` method. A destination that refuses the bytes during the flush, such as a full stick, is reported only there.
