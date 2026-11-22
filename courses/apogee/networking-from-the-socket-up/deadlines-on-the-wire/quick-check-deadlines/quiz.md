## When does the second `readFrame` give up?

- [ ] At t0+11s, 5 seconds after it starts
  > A deadline is an instant, not a span: nothing moves it when a new read starts.
- [x] At once: the deadline passed at t0+5s and still stands
  > The instant stays until something changes it, so every read after it fails right away.
- [ ] Never: the first read used the deadline up
  > A deadline isn't spent by one read; it holds for every read until it is changed.

## How does dishd remove the read deadline?

- [ ] `conn.SetReadDeadline(time.Now())`
  > That sets a deadline that has already passed, so every read fails at once.
- [ ] `conn.Close()`, then dial again
  > That removes the whole connection, not just its deadline.
- [x] `conn.SetReadDeadline(time.Time{})`, the zero time
  > The zero `time.Time` means no deadline: reads wait for as long as it takes again.

## A read has just failed with `os.ErrDeadlineExceeded`. What about the connection?

- [x] It works again once a later deadline is set
  > A deadline closes nothing: set a new instant and the next read waits for bytes again.
- [ ] It is closed, and dishd must dial again
  > Only the read gave up; the connection stays open.
- [ ] It works again at once, with nothing changed
  > The deadline that passed still stands, so every read fails until a new one is set.
