`now.Sub(b.Last)` is how long it has been since the bucket was last brought up to date. Only when that is above zero do you add tokens and move `b.Last` to `now`, so a clock that steps backward changes nothing.

The tokens to add are the elapsed time's `Seconds()` times `refillPerSec`; `min` keeps the sum at `bucketSize`. Then a request needs a whole token: with fewer than 1 left, return false without taking anything.

If a burst an hour later lets more than 10 through, you are not capping the tokens at `bucketSize`. If a clock that steps backward lets a request through, or lets two through a second after it, you added tokens or moved `b.Last` when `now` was before it.
