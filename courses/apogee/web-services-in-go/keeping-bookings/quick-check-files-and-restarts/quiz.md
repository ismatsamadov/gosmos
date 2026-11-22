## While `os.Rename` swaps the new file in, what can a program reading `bookings.json` see?

- [ ] Some bookings from the old file and some from the new one
  > The new bookings were written to the temporary file, never into `bookings.json` itself.
- [x] All of the old file, or all of the new one
  > The rename replaces the name in one step: a reader opens either the old file or the finished new one.
- [ ] An empty file, for a moment
  > The temporary file was written and closed before the rename, and the old file stays in place until the swap.

## At start, `bookings.json` holds half a JSON array. What should `OpenFileStore` do?

- [x] Return an error, so the service doesn't start
  > Someone must look at the file first; starting empty would let the next save write over every booking in it.
- [ ] Start with an empty store and carry on
  > The next booking's save would replace the whole file, and every booking in it would be lost.
- [ ] Keep the bookings that decoded and drop the rest
  > The broken half may hold bookings customers rely on; a store that guesses is worse than one that stops.

## The saved file holds `bk-4501` to `bk-4507`. Which ID does the first booking after a restart get?

- [ ] `bk-4501`, where every new store starts
  > That is what Monday's restart did, handing out IDs that customers already had.
- [ ] `bk-4507`, the last ID saved
  > `bk-4507` is already taken; the next ID is one more than the highest saved.
- [x] `bk-4508`, one more than the highest saved
  > `load` moves the count past the highest ID in the file, so no saved ID is handed out again.
