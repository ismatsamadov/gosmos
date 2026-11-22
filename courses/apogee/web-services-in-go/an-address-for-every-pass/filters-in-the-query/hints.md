Read `sat` and `station` from the query first. Check each one that isn't empty with `knownSat` or `knownStation`, and answer the 400 before the loop, quoting the value with `%q` as the messages do.

Start an empty list of `PassView`, not a nil one, and append each catalog pass whose satellite equals `sat`, or any satellite when `sat` is empty, and likewise for `station`.

If Brine-3 over Ystmark answers `null`, the list started as a nil slice. If `?sat=thistle` lists passes, the check is missing, or it runs only after the list is written.
