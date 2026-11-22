Wrap the body in `for { ... }`. Each turn, load `p.cur`; when there is a current schedule and `s.Version` is not greater than its `Version`, return `false`.

Otherwise call `p.cur.CompareAndSwap` with the schedule you loaded and `s`. When it returns `true`, return `true`; when it returns `false`, another offer changed `p.cur` after your load, so loop to load and compare again.

If version 1 replaces version 2, nothing compares the versions before the swap. If offers are still running after a second, a failed swap is tried again with the old load instead of a new one. If the test says `Offer` never calls `CompareAndSwap`, the check and a plain `Store` are still two steps another offer can come between.
