Name the result in the signature, `(err error)`, so that a deferred function can set it. Defer the function after the lookup and before the call to `run`: a deferred function only covers what happens after it is deferred.

In the deferred function, `recover()` returns `nil` when nothing panicked. Anything else means the command panicked, and calling `recover` has already stopped the panic.

Set `err` with `fmt.Errorf`, the text `internal error in %s` and the command's name, `words[0]`. Leave the rest of `dispatch` as it is, so unknown commands and ordinary errors don't change.
