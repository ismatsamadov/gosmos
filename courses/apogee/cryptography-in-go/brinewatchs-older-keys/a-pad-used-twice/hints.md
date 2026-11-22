`p.key[p.used:]` is the part of the pad nobody has used yet. Compare the message's length with that part's length, not the whole pad's, and put the number left in the error.

`subtle.XORBytes` writes as many bytes as the shorter of its two inputs, so XOR `msg` with the unused part. Only after that add `len(msg)` to `p.used`: a refused message returns before anything changes.

If `TestSomeoneKeptBoth` still reads January's secret, the second seal reused bytes the first had used. If `TestUsedOnce` says a call started in the wrong place, `p.used` grew by something other than the message's length.
