Make the group from `ctx` with `cfg.Slots` slots, and call its `Go` once for each station, in the loop the starter already has. The job's body is the starter's call to `planStation`, made with the job's own `ctx`, and it returns the error.

`Wait` returns the group's first error once every job has returned; return it before you build the schedule. Merging needs one sequence per station, and `slices.Values` makes one from a station's list.

If the test finds one station planning at a time, `Plan` still waits for each station before starting the next; if a thirteenth starts, the jobs aren't going through the group's slots.
