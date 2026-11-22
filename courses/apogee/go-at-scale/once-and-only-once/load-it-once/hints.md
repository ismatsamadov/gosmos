`t.once.Do` takes a function with no arguments and no results, and runs it only the first time any caller reaches it: the call to `loadNetwork` goes inside such a function.

In that function, store both of `loadNetwork`'s results in `t.net` and `t.err`. After `Do` returns, the work is done, by this caller or another, so return the two fields.

If the test says `loadNetwork` ran 100 times, it is still called outside `t.once.Do`. If callers get networks of their own, they return what their own call produced instead of the fields.
