Leave the current schedule alone: load it only to read its `Version`, and count 0 when there is none yet.

Make a new `&Schedule{...}` with that version plus one and `slices.Clone(contacts)` as its contacts, `Store` it in `p.cur`, and return it. Import `"slices"`.

If the schedule the Board held still changes, something writes to a schedule that was already published. If a later change to the caller's contacts shows up in the schedule, it keeps the caller's slice instead of a clone.
