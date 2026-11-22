## What does the random factor in each wait prevent?

- [ ] A receiver guessing when the next attempt will come
  > Nothing about the schedule is secret; the signature protects the event, not the timing.
- [ ] Waits that grow past five minutes
  > The cap does that; the factor only shortens a wait, to between half and all of it.
- [x] Deliveries that failed together retrying together
  > Without it, every event that failed in the same second would come back in the same second and knock the receiver over again.

## A receiver answers `404 Not Found`. What should the sender do?

- [x] Give up at once: asking again gets the same answer
  > A 4xx other than 408 and 429 says the request itself is wrong, and it won't be less wrong later.
- [ ] Try again on the schedule, as it would after a `503`
  > A 503 says the receiver is unwell for now; a 404 says there is nothing at that address to receive it.
- [ ] Try again once, at once, without waiting
  > The same request to the same address gets the same 404.

## Why must a receiver be ready to get the same event twice?

- [ ] The sender sends every event twice, to be safe
  > It sends once per attempt and stops at the first 2xx.
- [x] Its 2xx can be lost on the way back, and the sender then tries again
  > The sender can't tell a lost answer from an event that never arrived, so it retries; the event's ID lets the receiver drop the copy.
- [ ] Each retry carries a new event ID
  > Every attempt carries the same ID: that is how the receiver knows it is a repeat.
