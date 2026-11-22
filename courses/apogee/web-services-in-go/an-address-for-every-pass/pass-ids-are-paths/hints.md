Replace `{id}` in the pattern with three wildcards, one per segment, named as the task names them. The handler stays `s.getPass`.

In `getPass`, ask the request for each of the three values with `PathValue`, and join them in the pass ID's order, satellite, station, orbit, with a slash between each two.

If the 404 reads `no pass //` or is missing a part, a wildcard's name in the pattern differs from the one you ask `PathValue` for, which then returns an empty string. If every pass gets `404 page not found`, the pattern still has `{id}`.
