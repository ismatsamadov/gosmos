Three reads can meet the end of the stream: the size line's `readLine`, the chunk's `io.ReadFull`, and the trailer's `readHeader`. Every one of them is a cut, so none may return the body as if it were whole.

A small function that turns `io.EOF` and `io.ErrUnexpectedEOF` into `&RequestError{Status: 400, Msg: "recording cut off"}`, and returns any other error as it is, can wrap the error at each of those three returns.

If the test still finds `first piece ` kept, a branch returns the body with a nil error when `readLine` reports `io.EOF`: remove it. The trailer counts too, so check the error `readHeader` returns before returning the body.
