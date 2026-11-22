Write the counter with `binary.BigEndian.AppendUint64(nil, counter)` and pass those 8 bytes to the HMAC's `Write`; `Sum(nil)` then returns 20 bytes, SHA-1's size.

The offset is the last byte `& 0x0f`, a number from 0 to 15, so the 4 bytes from it always fit in 20. `binary.BigEndian.Uint32` reads them, `& 0x7fffffff` clears the top bit, and `% 1000000` keeps six digits.

If counter 0 gives `238872`, the top bit is still set; if counter 1 gives `717529`, the 4 bytes come from the start instead of the offset. If `081804` shows up as `81804`, format with `%06d`.
