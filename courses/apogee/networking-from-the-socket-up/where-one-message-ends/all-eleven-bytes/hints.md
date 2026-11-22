`io.ReadFull` takes the reader and a slice, and calls `Read` as many times as it takes to fill the slice. Make the slice `n` bytes long.

It returns the number of bytes it read and an error: `nil` once the slice is full, `io.ErrUnexpectedEOF` when the stream ended part way, `io.EOF` when it ended before the first byte. Return that error as it is.

If the PONG comes back as one byte, or three, a single `Read` is still deciding how much you get.
