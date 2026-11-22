---
slug: lost-or-just-refused
title: Lost or Just Refused?
type: code_tests
position: 3
difficulty: 8
topics: errors, networking
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Tuesday's log holds a second puzzle. Once the wind reached 72 km/h, Dish 2's controller held the dish stowed, its wind lock, and refused every command to unstow or move it with a NAK `wind lock`. dishd took every NAK for a dead line and redialed forty times, though the controller had answered each time.

A lost line means hang up and reconnect: the stream ended, a deadline passed, the connection was closed or reset, the controller said BYE, or its bytes can't be trusted. A NAK is an answer from a working line, and an ended context means stop. `errors.Is` finds an error through any wrapping, and `errors.AsType` finds one of a type:

```go
if nak, ok := errors.AsType[*NakError](err); ok {
	fmt.Println("refused, code", nak.Code)
}
```

## Your task

`lost` in `link.go` returns `err != nil`. Return true when `err` is or wraps one of these, and false for anything else:

```text
io.EOF  io.ErrUnexpectedEOF
os.ErrDeadlineExceeded
net.ErrClosed  io.ErrClosedPipe
syscall.ECONNRESET  syscall.EPIPE
ErrBadLength  ErrBadCRC  any *ByeError
```
