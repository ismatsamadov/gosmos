Start `safeHeader` with `h.Clone()`: the copy holds its own values, so you can change them without touching the request's header. Return the copy at the end.

For each of the four names, look up the copy's values under the name as written in the task, which is its canonical form, and set every one of them to `[redacted]`, so a header sent twice is still two values.

If the test says the request's own header changed, you redacted `h` instead of the copy. If a cookie added in lowercase slips through, look it up as `Cookie`: `Add` stored it under that name.
