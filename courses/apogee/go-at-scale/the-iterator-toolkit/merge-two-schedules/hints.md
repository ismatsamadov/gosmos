Pull from both sequences with `iter.Pull` and defer both stops. Then pull a first head from each, keeping its `ok` beside it.

Loop while either head is there. When `a` has a head and `b` has none, or `cmp` of the two heads is 0 or less, yield `a`'s head and pull its next; otherwise do the same with `b`'s.

If the test shows passplan 2's copy first, a tie went to `b`: equal heads must yield `a`'s. If the last values are missing, your loop ended as soon as one sequence did; if a loop that stops at one read every contact, both were collected first instead of pulled.
