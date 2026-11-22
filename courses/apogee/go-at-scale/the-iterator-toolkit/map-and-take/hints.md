`Map` is `Filter` without the test: range over `seq` and hand `f(v)` to `yield`, returning when it gives back false.

In `Take`, return before ranging when `n` is 0 or less, so `seq` is never asked. Otherwise count each value as it arrives, inside the function you return, and return as soon as `n` of them have gone to `yield`.

If the schedule was asked for 13 contacts, you checked the count only when the 13th arrived. If a second loop over the same `Take` gets a different number of values, the count lives outside the returned function, or you changed `n` itself.
