Make a slice of exactly `n` bytes and hand it to `io.ReadFull` with `br`, not the connection: the reader that read the head may already hold the body's first bytes.

`io.ReadFull` returns how many bytes it read and an error. `io.EOF` means none arrived and `io.ErrUnexpectedEOF` means some did; both mean the body ended early, and the count goes in the message.

If the server answers 408 in the last subtest, the body is still read line by line: the recording has no newline at its end, so the read waits for one until the deadline. Build the refusal as `&RequestError{Status: 400, Msg: ...}`, with the message from `fmt.Sprintf`.
