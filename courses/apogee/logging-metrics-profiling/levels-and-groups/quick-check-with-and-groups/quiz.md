## What does the last line write?

- [ ] `{"level":"INFO","msg":"pass started","program":"uplinkd","site":"LAN"}`
  > `With` doesn't change `l`: it returned a new logger, `u`, and only `u`'s records carry `program`.
- [x] `{"level":"INFO","msg":"pass started","site":"LAN"}`
  > `l` is the logger it always was; the attributes from `With` belong to `u` alone.
- [ ] Nothing, since `l` was used up when `With` made `u`
  > A logger is never used up: `With` leaves the old one working, unchanged, beside the new one.

## Through the text handler instead, how would the fourth line write `g`?

- [ ] `pass={sat_id=brine-3}`, an object as in JSON
  > Only the JSON handler nests a group; the text handler has no objects, so it puts the group's name in front of each key.
- [ ] `sat_id=brine-3`, the group's name dropped
  > The name is kept: without it, a `sat_id` in a group couldn't be told from one outside it.
- [x] `pass.sat_id=brine-3`, the group's name and a dot before the key
  > The text handler flattens a group into dotted keys, while the JSON handler writes it as a nested object.

## In the fourth line's JSON, what comes right after `"msg":"pass started"`?

- [x] `"program":"uplinkd"`, because attributes from `With` come before the record's own
  > The handler holds `With`'s attributes already written, and puts each record's own after them.
- [ ] `"pass":{"sat_id":"brine-3"}`, because a record's own attributes come first
  > The record's own attributes come after the ones `With` added, so `pass` follows `program`.
- [ ] `"pass":{"program":"uplinkd","sat_id":"brine-3"}`, since `With` joins the group
  > `With` added `program` outside any group; the group holds only what `slog.Group` was given.
