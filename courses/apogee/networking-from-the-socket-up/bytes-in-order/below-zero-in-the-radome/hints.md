`binary.BigEndian.Uint16` gives you the 16 bits as a `uint16`, from 0 to 65535. Wrap that call in `int16(...)` to read the same bits as a number from −32768 to 32767.

Then wrap the `int16` in `float64(...)` and divide by 10, so the division keeps its tenths.

If `ff 83` gives 6541.1, the bits were never read as signed; if it gives −12, the division happened before the conversion to `float64`.
