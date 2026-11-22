Add the field next to `commands` and set it in the struct `newShell` returns. `tonight()` returns a new slice each time, in the schedule file's order, so sorting is your command's job.

Sort a copy: clone `s.passes`, then call `slices.SortFunc` on the clone with a function that compares two passes' `AOS` fields with `Compare`. Sorting `s.passes` itself would change the shell's own order, and the next lesson checks that it doesn't.

Format both times with the layout `15:04` and the elevation with `%.1f`; everything else is single spaces, with the word `max` before the elevation. In `commandTable`, the new entry's `run` is `cmdPasses`.
