## After the loop has printed every ID, Nkem calls `next()` twice more. What do the calls return?

- [ ] The last ID and true, both times
  > `next` doesn't hand back the last value again; once the sequence has ended, it has nothing left to give.
- [x] `""` and false, both times
  > After the end, `next` returns the zero value and false, as often as it's called.
- [ ] `""` and false, then a panic
  > Calling `next` after the end is allowed, as often as you like, and changes nothing.

## A version of the loop breaks at its first Brinewatch ID, and an edit adds a second `stop()` after it, beside the deferred one. What happens?

- [ ] The second call panics, since the sequence is already stopped
  > `stop` may be called more than once; once the sequence has stopped, another call does nothing.
- [x] The first call ends the sequence and the second does nothing, and one of them had to run
  > A sequence left before its end waits inside its loop until `stop` ends it, and calling `stop` again is fine.
- [ ] Neither call was needed: breaking out of the loop stops the sequence
  > Breaking out of a loop that calls `next` doesn't touch the sequence, which waits inside its loop until `stop` runs.

## Can a second goroutine call `stop` while the first is inside a call to `next`?

- [ ] Yes: `stop` is safe from any goroutine at any time
  > `iter.Pull`'s two functions aren't safe to call from two goroutines at the same time, `stop` included.
- [ ] Yes, as long as each goroutine calls `stop` only once
  > How often doesn't matter; what's ruled out is two calls at the same time.
- [x] No: `next` and `stop` must never be called from two goroutines at once
  > `iter.Pull` allows one call at a time, and `Merge` makes all its calls from the one goroutine that ranges over it.
