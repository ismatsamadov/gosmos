Turn each count into a `float64` before dividing, as `availability` does, and the target's share the same way: Lantern Atoll's 9900 leaves 100 basis points of 10,000, a share of 0.01.

Return 0 first when `booked` or `lost` is 0. With no loss the rate is 0 whatever the target, and that check also keeps 0/0, which is NaN, out of the result.

If the test shows NaN, a zero count reached the division. A loss against 10000 needs no special case: dividing a positive share by 0.0 already gives `+Inf`.
