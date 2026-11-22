The check goes right after `n` is known and before the `make`: once the payload is being read, a bad length has already done its harm.

Both ends matter: lengths from 9 to 1024 are frames, and everything below or above is refused. `fmt.Errorf` with `%w` keeps `ErrBadLength` inside the error you return.

If the test says `readPayload` read on past the length, the check comes too late or misses one side; if the message is wrong, it should read `bad length` and the number, one space between.
