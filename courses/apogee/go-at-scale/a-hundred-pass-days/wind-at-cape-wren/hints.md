In `heard`, start from `slices.Clone(stows)`. With `slices.IndexFunc`, find a stow at `st.Station` whose `To` is zero.

If there is one and `st.To` isn't zero, set that stow's `To` and return the list; otherwise append `st`. In `Run`, at the TODO, hand an event's stow, when it has one, to `heard` and keep the result in `r.stows`, before the line that reads `r.stows` for the run.

If the end leaves two stows, the open one wasn't found: it must be at the same station and have a zero `To`.
