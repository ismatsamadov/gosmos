## With no constructor, is `up.mu` ready to lock?

- [x] Yes: a zero `sync.Mutex` is an unlocked mutex, ready to use
  > The zero value of `sync.Mutex` is a working, unlocked mutex, so a struct that holds one needs no setup.
- [ ] No: `mu` is nil until you assign it `&sync.Mutex{}`
  > `mu` is a `sync.Mutex` value, not a pointer, so it is never nil, and its zero value already works.
- [ ] No: it has to be unlocked once before its first `Lock`
  > A new mutex already starts unlocked, and unlocking one that isn't locked is a fatal error.

## `Send` returns `ErrUnsigned` early. What happens to `u.mu`?

- [ ] It stays locked, because the early return skips `Unlock`
  > Without `defer` it would; a deferred call runs whenever `Send` returns, early or not.
- [x] It is unlocked: the deferred `Unlock` runs on every return
  > `defer` schedules `Unlock` right after `Lock`, so every way out of `Send` releases the mutex.
- [ ] It is unlocked only when `Send` returns `nil`
  > A deferred call doesn't look at the result: it runs on every return, with or without an error.
