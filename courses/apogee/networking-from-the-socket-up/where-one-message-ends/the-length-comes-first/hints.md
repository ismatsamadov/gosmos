Read the length into a two-byte array or slice with `io.ReadFull(br, ...)`, then pass those two bytes to `binary.BigEndian.Uint16`.

Make a slice as long as that number and fill it with `io.ReadFull` as well. Return each `ReadFull` error as it is: `io.EOF` before the length means the stream ended between frames, `io.ErrUnexpectedEOF` inside one.

If the WELCOME comes back cut short after its `0x0a`, something is still reading up to a newline.
