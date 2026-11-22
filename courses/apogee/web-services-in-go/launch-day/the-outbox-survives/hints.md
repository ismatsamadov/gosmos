Marshal an `outboxFile` whose `NextEvent` is `o.next` and whose `Pending` is `o.pending`, then hand the bytes to `replaceFile` with `o.path`, which puts them in place whole; return the first error.

In `OpenOutbox`, read the file with `os.ReadFile`. `errors.Is(err, fs.ErrNotExist)` tells a missing file, which is no error, from one that can't be read, which is.

Decode into an `outboxFile`, then copy its `NextEvent` into `next` and its `Pending` into `pending`. If the reopened outbox hands out `evt-000001` again, `next` never came from the file.
