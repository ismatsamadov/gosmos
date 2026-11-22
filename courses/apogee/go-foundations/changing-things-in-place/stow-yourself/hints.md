Each method works on a copy of `yst` and changes only that copy. What kind of receiver would let it reach `yst` itself?

In both methods, change the receiver from a `Station` to a pointer to a `Station`. The bodies and the calls stay as they are.
