## Why does each wait carry a random factor from 0.5 to 1?

- [ ] So that no wait is ever shorter than half a second
  > A fixed half would do that too; the randomness is there so that waits differ from one link to the next.
- [x] So that Dish 1's and Dish 2's links, dropped by the same storm, don't redial in step
  > Two links that fail together would otherwise dial at the same instants every time; the jitter spreads their dials apart.
- [ ] So that the controller can tell a redial from a first connection
  > The controller sees a new connection either way; nothing in a dial carries the wait that came before it.

## Why do the waits stop growing at 30 seconds?

- [ ] Because the controller hangs up on a client that is silent for 30 seconds
  > That idle limit counts on an open connection; between dials there is no connection to keep alive.
- [ ] So that the count of failed dials can never overflow
  > The count goes on growing past the cap; only the wait stops at 30 s, and that number comes from how soon dishd must be back.
- [x] So that after a fourteen-minute outage dishd dials again within half a minute of the line's return
  > Uncapped, the waits would keep doubling, to eight minutes and more, and dishd could sleep through most of a pass after the line was back.

## A wait inside `redial` returns `context.Canceled`. What does `redial` do?

- [x] Return that error at once, since the caller has asked it to stop
  > A canceled context means stop, not try again, and returning its error tells the caller exactly that.
- [ ] Dial once more, in case the line has come back
  > Once the context has ended, nothing new should start: dishd is shutting down, or the caller has given up.
- [ ] Wait 30 seconds and start over from the first delay
  > The wait ended because the context did; waiting longer ignores the one signal that says stop.
