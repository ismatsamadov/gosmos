`Header.Get` works on a nil header too, returning `""`, so ask the trailer first and fall back to the header only when it has no value. When both are empty, return the CRC of the body with no error.

`strconv.ParseUint(sent, 16, 32)` reads the digits in either case; check also that there are exactly 8 of them. Then compare `uint32` of the result with `crc32.ChecksumIEEE(body)`.

If a chunked upload with a wrong trailer is still kept, the CRC is read from the header alone. The mismatch message prints both values with `%08x`, so a sender's capitals come back in lowercase.
