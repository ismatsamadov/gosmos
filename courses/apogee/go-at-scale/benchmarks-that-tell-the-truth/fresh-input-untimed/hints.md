Inside the loop, before the call to `planDay`, stop the timer, make the copy with `slices.Clone(day)`, and start the timer again. Add `slices` to the imports.

Hand `planDay` the clone, not `day`: `day` itself must stay as `stationDay` predicted it, so that every clone starts unsorted.

If the allocations per turn are one too many, the clone is made while the timer runs. If they are 0, the timer is still stopped when `planDay` runs.
