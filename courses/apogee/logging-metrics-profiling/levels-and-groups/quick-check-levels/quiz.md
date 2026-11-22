## Lantern Atoll's handler is at WARN, which is 4. What does it do with an INFO record, which is 0?

- [ ] It writes it, since INFO is the level every program starts at
  > Programs start at INFO, but this handler's minimum is WARN, and it decides by its own minimum.
- [x] It drops it, since 0 is below 4
  > A handler writes a record only when the record's level is at least its own: here WARN and ERROR get through, INFO and DEBUG don't.
- [ ] It writes it at WARN, raising it to the handler's level
  > A handler never changes a record's level; it only decides whether to write the record.

## The night shift sets the `LevelVar` to DEBUG while two passes are logging. What happens?

- [x] Every record checked after the change is checked against DEBUG, with no restart
  > A `LevelVar` is safe to change while other goroutines log, and the handler reads it again for every record.
- [ ] Nothing until uplinkd restarts, since the handler copied the level when it was made
  > The handler holds the `LevelVar` itself, not a copy of its value, so the next record already sees the change.
- [ ] The passes must stop logging first, or a record could be written half at each level
  > Each record is checked once, as a whole, and a `LevelVar` changes in one step that no goroutine can see half done.
