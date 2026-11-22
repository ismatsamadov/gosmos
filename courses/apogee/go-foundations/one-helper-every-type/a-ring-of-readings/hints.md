`NewRing` has to remember how many items to keep, and `Push` behaves differently before and after the ring fills up. How can `Push` tell which case it's in?

Before the ring is full, `append` to `items`. After that, write over the item at index `next`, then move `next` one place on, wrapping back to 0 at the end with the `%` operator.

Once the ring is full, the oldest item sits at `next`. For `Items`, make a new slice and append two parts of `items` to it: first from `next` to the end, then everything before `next`.
