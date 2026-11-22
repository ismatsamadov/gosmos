Collect the station's passes into a new slice, then sort it with `slices.SortFunc`, comparing `AOS` values with their `Compare` method. Other stations' passes never matter.

Keep a candidate start, beginning at `from`. Walk the sorted passes: if the switch ends no later than the pass's AOS minus 5 minutes, it fits, and you can stop; if the candidate is before the pass's LOS plus 5 minutes, move it there.

When the candidate's switch would end after `from` plus 24 hours, return the error with that instant in RFC 3339. If PEL comes out at 01:00:00, the passes weren't sorted, or the walk skipped Pipit-3's pass, which is already under way at 01:00.
