Each call returns the slice with the new bytes added: assign it back to `dst` every time, and keep what `dst` held before the call in front of the new line.

Eight appends make the line: the AOS, a space, the LOS, a space, the dish, a space and the pass ID, then the newline. A string's bytes go in with `append(dst, s...)`, one byte with `append(dst, ' ')`.

`strconv.AppendInt` takes an `int64`, so convert the dish. If the times come out without their leading zero, the layout isn't `time.TimeOnly`.
