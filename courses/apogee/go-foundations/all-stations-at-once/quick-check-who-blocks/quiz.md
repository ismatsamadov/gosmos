## With nobody receiving, what does `a <- "PEL"` do?

- [ ] It stores `"PEL"` in the channel and moves on
  > Storing a value needs a buffer, and `a` has none, so there is nowhere to put it.
- [x] It waits until another goroutine receives from `a`
  > An unbuffered send hands its value straight to a receiver, so it waits for one to arrive.
- [ ] It drops the value, since nobody is listening
  > A channel never drops a value: the send waits until it can deliver it.

## With nobody receiving, what does `b <- "SQT"` do?

- [x] It puts `"SQT"` in the buffer and moves on
  > `b` has room for two values, so the send completes at once; only a third send with no receive in between would wait.
- [ ] It waits until another goroutine receives from `b`
  > That is how an unbuffered channel behaves; `b`'s buffer takes the value while it has room.
- [ ] It panics, because nobody receives
  > Having no receiver never makes a send panic: it either waits or, with room in the buffer, completes.

## `main` returns while another goroutine is still waiting to send on `a`. What happens?

- [ ] `main` waits for that goroutine to finish first
  > `main` never waits for other goroutines by itself; that is what `wg.Wait()` is for.
- [ ] The program stops with a deadlock error
  > A deadlock needs every goroutine stuck, `main` included; here `main` has returned, and the whole program has ended with it.
- [x] The program ends, and the waiting goroutine ends with it
  > When `main` returns, the whole program exits, whatever its other goroutines are doing.
