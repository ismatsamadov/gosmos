Add `cmp`, `os`, `path/filepath`, `slices` and `strings` to the imports. Then sort `files` with `slices.SortFunc`: compare the days with `strings.Compare` and, for the same day, the numbers with `cmp.Compare`, joined with `cmp.Or`.

Add up the files' sizes, then work out the oldest day to keep: `now` in UTC, `keepDays` days back with `AddDate`, written with `time.DateOnly`. Walk the sorted files and stop at the first one that is both that day or later and in a total that fits; delete the others with `os.Remove`, taking each one's size off the total.

If `.10` goes before `.9`, the files were left in the order of their names: compare `N`, not `Name`. If 20 November goes, either a day exactly 14 days back counts as too old in your test, where only an earlier day goes, or `now` wasn't turned to UTC first: in its own zone it is already the 5th. If the current file goes, you listed the folder yourself: `rotatedFiles` never lists it.
