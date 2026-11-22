Take the passes from the cache and sort a copy by AOS, as `cmdPasses` does. The schedule lists them station by station, so the first pass after now in that order is not always the next one.

With the time from `s.cache.now()`, the first sorted pass whose `AOS.After(now)` is true is the answer. A pass whose AOS is exactly now has already begun, so it doesn't count.

Print the satellite, the station, the AOS through the `15:04` layout and `p.AOS.Sub(now)` through `%v`, in the format of the first example line. When the loop ends without a match, print the line for no passes left.
