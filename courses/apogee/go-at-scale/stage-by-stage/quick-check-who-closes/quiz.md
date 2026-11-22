## Which goroutine closes `frames`?

- [ ] The one in `newest`, once it has received the last frame
  > A receiver never closes: it can't know whether more is coming, and its close would make the sender's next send panic.
- [x] The one `parse` starts, after it has sent its last frame
  > The stage that sends on a channel owns it, so only it knows when nothing is left to send.
- [ ] The one running `main`, after `newest` returns
  > `main` never sends on `frames`, and by then `parse` has closed it already.

## What does adding `close(lines)` to this code do?

- [x] The build fails: `lines` is a `<-chan Heard`, which can only receive
  > A receive-only channel can't be closed, so the compiler refuses the call: only `read`, which made the channel, can close it.
- [ ] `read`'s next send panics
  > That would happen with a channel `main` could close; `read` hands out a receive-only one, so the program never builds.
- [ ] `parse`'s loop over `lines` ends early
  > The program never runs: the compiler refuses to close a receive-only channel.

## When does the loop in `newest` end?

- [ ] The moment `parse` closes `frames`, even with frames still in it
  > A `range` receives every value still in the channel before it ends.
- [ ] When `frames` is empty for a moment
  > An empty channel makes the loop wait for the next value; only a close ends it.
- [x] Once `frames` is closed and every frame in it has been received
  > A `range` over a channel drains what is left after the close, then ends.
