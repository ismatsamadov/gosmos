Look at the error first. `errors.Is(err, bufio.ErrBufferFull)` means the line filled the buffer, so return `errLineTooLong`. `io.EOF` with some bytes already read means the stream ended inside a line.

A good line ends in `\r\n` and holds no other `\r`, which some programs would read as the end of a line; anything else that ended in `'\n'` is the 400. An empty stream stays `io.EOF`: that is a client that said nothing.

Turn the line into a string before the next read reuses the buffer, and cut its CR LF off. If the 64 KiB test says more than 1024 bytes were read, `ReadString` is still doing the reading.
