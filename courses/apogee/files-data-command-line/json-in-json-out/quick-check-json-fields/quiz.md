## Which key does json.Marshal write for SatID?

- [ ] "SatID"
  > The tag replaces the Go name on the wire.
- [x] "sat_id"
  > The tag names the key for both encoding and decoding.
- [ ] "satid"
  > Encoding never changes a name on its own; the key comes from the tag.

## What does json.Marshal write for station?

- [ ] "station":"SQT"
  > A lowercase field is unexported, and the json package cannot see it.
- [ ] "Station":"SQT"
  > Marshal does not capitalize names; it skips fields it cannot see.
- [x] Nothing at all
  > Unexported fields are left out of the JSON, and Unmarshal never fills them either.

## What does json.Unmarshal return for {"id":"bk-2041", cut off there?

- [x] A *json.SyntaxError, with the booking unchanged
  > Unmarshal checks that the whole input is JSON before it fills any field.
- [ ] nil, with ID set to bk-2041
  > Broken input is an error, not a partial success.
- [ ] nil, with every field left empty
  > An empty booking with no error would hide the damage.
