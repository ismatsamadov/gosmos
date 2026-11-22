## What happens when the first snippet runs?

- [ ] It prints `ACK 118`
  > The send waits for a receiver, and the only receive comes after it, so `Println` is never reached.
- [ ] It waits forever without a message
  > Here no goroutine can ever move again, and Go notices that and stops the program instead of letting it hang.
- [x] It stops with `fatal error: all goroutines are asleep - deadlock!`
  > `main` waits on its own send, and no other goroutine exists to receive, so nothing can ever move again.

## What does the second snippet do?

- [x] The second `close` panics with `close of closed channel`
  > A channel can be closed once; closing it again is a bug, and Go reports it with a panic.
- [ ] The second `close` is ignored, since the channel is already closed
  > Go doesn't overlook a second close: it panics, so the bug shows up at once.
- [ ] The compiler rejects the second `close`
  > The compiler doesn't track whether a channel is closed; the check happens while the program runs.

## Once `frames` is closed and empty, what does `seq := <-frames` do?

- [ ] It waits for the next value to be sent
  > Nothing more can be sent on a closed channel, so the receive doesn't wait.
- [x] It returns 0 at once
  > A closed, empty channel gives its element type's zero value right away; `seq, ok := <-frames` would also set `ok` to `false`.
- [ ] It panics
  > Receiving from a closed channel is safe; sending on it or closing it again is what panics.
