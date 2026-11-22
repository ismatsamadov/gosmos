Every field is read through `order`, so the order is chosen in one place: the line that sets it.

`binary.BigEndian` and `binary.LittleEndian` are both a `binary.ByteOrder`, so either can be stored in `order`.

If the azimuth still prints 382.15, the bytes `47 95` are still read low byte first.
