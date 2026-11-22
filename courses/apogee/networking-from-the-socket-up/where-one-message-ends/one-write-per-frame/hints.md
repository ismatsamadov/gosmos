`encodeFrame(f)` already returns the whole frame, length and CRC included, so there is nothing left to assemble.

Call `Write` once with all of it. It returns a count and an error; return the error.

If the test counts two `Write`s for one frame, the length still goes out on its own. If it says the closed line's error is missing, the `Write`'s error is dropped instead of returned.
