`TeeReader` wraps the reader; it is not another destination for `io.Copy` here. Feed the wrapped reader into `ReadAll`.

Keep the read error, since it can come from either the source or the log writer.
