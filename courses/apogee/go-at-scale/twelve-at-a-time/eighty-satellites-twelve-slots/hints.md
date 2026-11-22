Start a group with `newGroup`, giving it `ctx` and `slots` as its limit, and make the map. For each satellite, hand the group a job that calls `predict` with the job's own context and returns its error.

When a prediction succeeds, lock a `sync.Mutex`, store the passes under the satellite's ID, and unlock. After the loop, `g.Wait()` returns the first failure: return `nil` and that error, or the map and `nil`.

If predictions still start after brine-17 fails, a job runs outside the group. If the running ones never see their context end, they are given the context `predictAll` received instead of the group's.
