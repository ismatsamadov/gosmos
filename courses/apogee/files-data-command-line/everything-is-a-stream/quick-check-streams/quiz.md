## Which interface describes a source that supplies bytes?

- [x] io.Reader
  > Its Read method supplies bytes to the caller's buffer.
- [ ] io.Writer
  > A writer accepts bytes from its caller.
- [ ] Only os.File can be a source
  > In-memory readers also satisfy io.Reader.

## Which argument of io.Copy receives the bytes?

- [ ] The second argument
  > The second argument is the reader supplying the bytes.
- [x] The first argument
  > Copy writes the source's bytes into the destination supplied first.
- [ ] Both arguments are modified into strings
  > Copy transfers bytes through interfaces rather than converting its arguments.

## Copy wrote four bytes before a read error. What must the helper return?

- [ ] Zero bytes and nil
  > That hides both the partial transfer and the failure.
- [ ] Four bytes and nil
  > The count does not make the failure disappear.
- [x] Four bytes and the error
  > Both results describe what happened before copying stopped.
