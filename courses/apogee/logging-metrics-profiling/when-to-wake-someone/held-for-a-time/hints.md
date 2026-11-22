`beforeFiring` has three jobs, in this order: forget the wait when `s.Value` is below `r.FireAt`; start the wait when `st.since` is still zero; and fire only when the wait has lasted `r.For`.

`st.since.IsZero()` tells you the wait hasn't started, `time.Time{}` sets it back to zero, and `s.At.Sub(st.since)` is how long the value has stayed at or above `FireAt`. Keep the last two lines, which fire the rule.

If the dip at 09:10 still fires at 09:15, a sample below `FireAt` doesn't clear `st.since`. If 15 minutes exactly don't fire, the comparison with `r.For` leaves out the boundary: at least `For` is enough.
