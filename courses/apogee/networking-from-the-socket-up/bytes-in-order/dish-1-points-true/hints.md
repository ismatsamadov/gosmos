After `send`, call `c.recv()` once. Its error, a BYE as a `*ByeError` or a damaged frame as `ErrBadCRC`, goes back to the caller unchanged.

Then look at the answer: `TypeAck` with the command's own `Seq` means done, and `TypeNak` becomes `nakError(ans)`, which carries the code and the reason.

Everything else, an ACK of another seq included, is an error made with `fmt.Errorf`, the type as `%#02x` and both sequence numbers as plain numbers. If ACK 118 came back as nil, the seq was never compared.
