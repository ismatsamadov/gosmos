A random nonce doesn't depend on the nonces before it, so `Noncer` needs no field at all: delete `next`, and the lines that write it into the nonce and increase it.

Keep the 12-byte slice and the `return`, and fill the slice with `rand.Read` between them. `encoding/binary` then has nothing left to do, so remove its import, and the `_ = rand.Read` line with it.

If `TestRestart` still finds two equal nonces, some bytes still come from the count: all 12 must come from `rand.Read`.
