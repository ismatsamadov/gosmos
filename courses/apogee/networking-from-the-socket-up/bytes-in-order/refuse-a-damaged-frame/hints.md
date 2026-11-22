Check the length first: under 9 bytes there is no type, seq and CRC to split, and slicing would fail. Wrap `ErrBadLength` with `%w` and the length, as `bad length 8`.

The CRC is the payload's last four bytes, read with `binary.BigEndian.Uint32`. Compute `crc32.ChecksumIEEE` over every byte before them and compare the two numbers.

Only when they match, split the fields. If the test says a damaged frame came back with its fields filled in, the fields were decoded before the CRC was checked.
