Keep a count before the loop, add `len` of the contacts `assign` returns in every turn, and after the loop call `b.ReportMetric` with the count divided by `b.N`, both as `float64`, and the unit `"contacts/op"`.

For step 2, pass `slices.Values(day)` where `stream` is now, then delete `next`, `stop` and `stream`, and the `iter` import with them.

3.8 is 76 contacts over 20 turns: only the first turn planned. 1520 is all 20 turns' contacts, not yet divided by `b.N`.
