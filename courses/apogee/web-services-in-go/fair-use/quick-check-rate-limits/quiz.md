## What does a bucket that holds 10 tokens promise a customer?

- [ ] 10 requests a second, all day long
  > The bucket refills at 1 token a second; 10 is how many it can hold at once.
- [x] Up to 10 requests at once after a quiet spell, then 1 a second
  > A full bucket pays for a burst of 10; after that, each second adds one token back.
- [ ] 10 requests in all, until the next day
  > The bucket refills all the time, at 1 token a second, never in one go at midnight.

## Brinewatch's bucket emptied at 09:00:00. Of a burst at 09:00:02.5, how many are answered?

- [x] 2
  > Two and a half seconds add 2.5 tokens, and each request needs a whole one.
- [ ] 3
  > Rounding 2.5 up would spend half a token that isn't there yet.
- [ ] 10
  > The bucket fills at 1 token a second, so it needs 10 quiet seconds to be full again.

## Why do the bucket's tests set a fake clock instead of sleeping?

- [ ] A fake clock runs faster than the real one
  > A fake clock doesn't run at all: the test sets its time.
- [ ] The service can't read the real clock
  > The real service reads the real clock through the same `now` function the tests replace.
- [x] The test sets each time exactly, so the answers are exact and nothing waits
  > Stepping the clock 2.5 seconds ahead takes no time, and the result can't depend on how busy the Mac is.
