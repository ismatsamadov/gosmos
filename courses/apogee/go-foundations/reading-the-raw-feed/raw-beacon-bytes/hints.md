Three things can be wrong with a beacon. In which order must you check them, so that each check can rely on the ones before it?

`hex.DecodeString` gives you the bytes and an error; after that, `len(b)` and `string(b[:2])` tell you whether you are holding a beacon at all.

`binary.BigEndian.Uint16` reads the two bytes of the slice you give it: `b[3:5]` for the sequence number, and the matching slice for bytes 6 and 7.
