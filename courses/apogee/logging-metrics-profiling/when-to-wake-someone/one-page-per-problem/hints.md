Two things change in `whileFiring`. In the branch that resolves, reset the state as before, then return `""` unless `r.Severity` is `Critical`.

After that branch the rule is still firing. Stay quiet when `r.Repeat` is 0, or when less than `r.Repeat` has passed since `st.sent`; otherwise remember this sample's time in `st.sent` and send a FIRING.

If the reminder comes a minute late, or never, check the comparison: a sample exactly `Repeat` after the last FIRING sends the next one. If a reminder comes every minute after the first, `st.sent` isn't updated when one is sent.
