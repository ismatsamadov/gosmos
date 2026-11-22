Look at `h["content-length"]`, not `h.Get`: the slice tells one field from two. No field at all is 0 for GET and HEAD and the 411 for PUT.

`strings.Trim(v, "0123456789")` leaves something behind whenever `v` holds anything but digits, and an empty value is bad too. Quote the value with `%q` in the 400's message.

`strconv.ParseInt(v, 10, 64)` fails only for a number too big for an `int64`, which is the 413 as well; compare with `max` only when it is above 0. If `+5` comes back as 5, the value reached `strconv.Atoi` before its digits were checked.
