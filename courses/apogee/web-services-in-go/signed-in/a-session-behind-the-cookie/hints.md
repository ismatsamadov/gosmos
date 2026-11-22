In `Start`, make the ID first, then lock `s.mu` and keep, under that ID in `s.byID`, a `session` for the customer whose `Ends` is `sessionLife` after `s.now()`.

In `Lookup`, a session has ended when the time is not before its `Ends`, the very instant included, so ask `Before` and negate the answer. The built-in `delete` removes an ended session from the map.

A struct you read from a map is a copy. After moving its `Ends`, store it back under its ID, or the renewal is lost and a busy session still ends 30 minutes after it started.
