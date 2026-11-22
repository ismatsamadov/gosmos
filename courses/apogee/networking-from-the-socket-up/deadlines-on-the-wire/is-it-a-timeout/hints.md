`errors.Is(err, target)` looks through every layer of wrapping for `target`. Here the targets are `os.ErrDeadlineExceeded`, `io.EOF` and `io.ErrUnexpectedEOF`.

`timeout.String()` writes the duration as `50ms` or `5s`, so the deadline's line is the name, `: no answer in ` and that text.

If the iotest row says `no answer`, the check reads the error's text; if the context row does, it asks the error's `Timeout` method or looks for `context.DeadlineExceeded`. Only `errors.Is` with `os.ErrDeadlineExceeded` names the line's own deadline.
