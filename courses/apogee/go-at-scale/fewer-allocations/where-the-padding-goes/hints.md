Each field needs a line like the one for `Sat`: that field's offset, that field's size and its own number, from `'2'` for `Station` to `'7'` for `Booked`.

`unsafe.Offsetof` takes a field selected from a value, such as `g.Station`, and reports where that field starts inside `g`; `unsafe.Sizeof` of the same field is how many bytes it takes.

If dots sit where a field should be, compare that field's line with `Sat`'s: its offset, its size and its number must all belong to the same field.
