By the time of the `TODO`, `b` holds the length and then the type, seq and body. The bytes the CRC covers are `b` without its first two.

Pass that slice to `crc32.ChecksumIEEE` and append the result with `binary.BigEndian.AppendUint32`, in place of the zero.

If a PING ends `c0 a3 dd d6`, the CRC included the length; if it ends `7d a4 e2 bc`, the right number went out little-endian.
