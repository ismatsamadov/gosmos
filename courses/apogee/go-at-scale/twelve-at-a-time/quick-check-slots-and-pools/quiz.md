## Twelve checks hold their slots. What does the thirteenth `Acquire` do?

- [ ] Returns an error, since every slot is taken
  > `Acquire` returns an error only when its context ends; while the run goes on, the send into a full channel waits.
- [x] Waits until one of the twelve calls `Release`
  > The send into a full channel waits for room, and a `Release` makes room by receiving one value.
- [ ] Takes a thirteenth slot, and the channel grows
  > A channel's room is fixed when `make` creates it: twelve slots stay twelve.

## A check fails, so its goroutine returns without calling `Release`. What happens to its slot?

- [ ] Go frees it when the goroutine ends
  > Nothing ties the value in the channel to the goroutine that sent it: only a receive takes it out.
- [ ] The next `Acquire` takes it over
  > The value is still in the channel, so the next `Acquire` finds one slot fewer, not a free one.
- [x] It is lost for good: after twelve failures, every `Acquire` waits until the run is canceled
  > So give the slot back on every way out: `defer s.Release()` as the goroutine's first line.

## For eighty jobs, how many goroutines do twelve slots start, and how many does a pool of twelve workers?

- [x] Slots start eighty, twelve working at a time; the pool starts twelve
  > A semaphore starts a goroutine for every job and lets twelve hold slots, while a pool's twelve goroutines take job after job.
- [ ] Twelve each
  > Slots start a goroutine for every job; only the number working at once is twelve.
- [ ] Eighty each
  > A pool starts its workers once, whatever the number of jobs, and each takes job after job.
