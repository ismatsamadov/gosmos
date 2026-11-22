## Besides a connection, what is `conn`?

- [ ] Only an `io.Reader`: writing needs a second value
  > A `net.Conn` has a `Write` method too, so one value carries both directions.
- [x] Both an `io.Reader` and an `io.Writer`
  > It has `Read` and `Write` with exactly their signatures, so `io.Copy`, `fmt.Fprintf` and every other reader or writer function works with it.
- [ ] Neither: a connection only has methods of its own
  > Its `Read` and `Write` match `io.Reader` and `io.Writer` exactly, so it satisfies both interfaces.

## The controller has sent its 21-byte banner. What can `n` be after the `Read`?

- [x] Any number from 1 to 21
  > `Read` returns the bytes that have arrived so far, at least one; the rest may still be on the way.
- [ ] Always 64, the size of `buf`
  > The buffer's size is the most `Read` can return, and it never waits to fill it.
- [ ] Always 21, because `Read` waits for the whole line
  > `Read` knows nothing about lines: it returns what has arrived, which can be part of the banner.

## Why does `defer conn.Close()` come after the error check?

- [ ] So that the connection closes before the `Read` runs
  > A deferred call runs when the function returns, after the `Read`, wherever the `defer` stands.
- [ ] Because Go allows a `defer` only after an `if`
  > A `defer` can stand anywhere in a function; the order here is about the value of `conn`.
- [x] A failed dial leaves `conn` nil, and deferring its `Close` would panic
  > The `defer` statement evaluates `conn.Close` at once, and on a nil connection that panics before anything returns.
