## Why does `readFrame` read from `br`, the reader the banner came through, and not from a new one?

- [ ] A new reader would read the banner a second time
  > The banner is already off the connection; what a new reader misses is anything `br` read past it.
- [x] `br` may already hold the first frame, read along with the banner
  > A `bufio.Reader` reads ahead, so a WELCOME sent in the banner's own write sits in `br`'s buffer, where a new reader would never find it.
- [ ] A connection can have only one reader wrapped around it
  > Any number of readers can wrap a connection; the trouble is what the first one has already taken off it.

## Why is the deadline moved before every read, not set once when the call opens?

- [x] Set once, it would end every call 90 s after it opened, with frames still coming
  > Every pass lasts at least two minutes; moved before each read, the deadline fires only after 90 s without a byte, as on a dead fiber.
- [ ] A deadline set once stops counting after the first read
  > A deadline holds for every read until it is moved or cleared, which is the trouble with setting it once.
- [ ] Moving it is what ends the call at the pass's LOS
  > The pass's context closes the call at LOS; the deadline only watches for silence.

## A frame's CRC doesn't match, and a later frame's length reads 8. What does each mean?

- [ ] Both mean the line is lost: hang up and dial again
  > A bad CRC spoils one frame whose length still showed where it ends, so the next frame starts in the right place.
- [ ] The CRC means the line is lost; the length means drop that frame and read on
  > It is the other way round: a wrong CRC concerns only its own frame, and a bad length leaves no way to find the next.
- [x] The CRC means drop that frame and read on; the length means the line is lost
  > A length under 9 can't be a frame, so nothing after it starts where it should, while a wrong CRC damages only its own frame.
