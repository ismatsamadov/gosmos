## What does the first `Println` show?

- [ ] The real time at which the test started
  > A bubble's clock is fake from the start: `time.Now()` inside it never reads the Mac's clock.
- [x] `2000-01-01 00:00:00 +0000 UTC`
  > Every bubble's clock starts at midnight UTC on 1 January 2000, which is why tests shift it to the day they need.
- [ ] `0001-01-01 00:00:00 +0000 UTC`, the zero time
  > That is what an unset `time.Time` holds; a bubble's clock starts at a real instant, the start of 2000.

## If the goroutine read from a real network connection instead of `quit`, what would the `Sleep` do?

- [ ] Return at once, ten fake minutes later
  > The clock jumps only when every goroutine is blocked on something only the bubble can wake, and a real read can be woken from outside it.
- [ ] Panic, since a bubble can't use the network
  > Nothing stops the read; it only keeps the clock still, which is why tests in a bubble use in-memory fakes such as `net.Pipe`.
- [x] Wait until the read is over: until then the clock can't move
  > A goroutine in a real read isn't blocked for good as far as the bubble knows, so the bubble isn't idle and its clock stays put.

## What happens when the function returns while the goroutine still waits on `quit`?

- [x] The bubble reports a deadlock: the test fails with a panic, and no later test runs
  > The goroutine is blocked on a channel only the bubble could send on, and the function that could have done it has returned.
- [ ] The goroutine is left behind, and the test passes
  > `synctest.Test` waits for every goroutine in its bubble, and one that can never be woken is a deadlock.
- [ ] `synctest.Test` waits for it forever
  > Once the function returns, nothing in the bubble can ever send on `quit`, so the bubble knows the wait can't end and says so.
