---
slug: nobody-listening
title: Nobody Listening
type: code_tests
position: 6
difficulty: 6
topics: errors, networking
created: 2026-11-07T07:44:27Z
modified: 2026-11-07T07:44:27Z
---
Dish 2's controller is still in its crate, and every time dishd tries it, its log says `connect: connection refused`. Sunniva wants one plain line instead: `YST-2: not listening`.

When nothing listens on a port, the Mac refuses the connection at once. `Dial`'s error then wraps the system's own error, `syscall.ECONNREFUSED`, inside layers that name the address. `errors.Is` looks through every layer for the error you name:

```go
_, err := conn.Read(buf)
if errors.Is(err, syscall.ECONNRESET) {
	fmt.Println("the other end hung up hard")
}
```

## Your task

In `reach`, return the name and `not listening`, as in `YST-2: not listening`, when the dial was refused. Keep the error's text for every other failure.

> [!TIP]
> Don't search the message for "refused": its text names an address that changes, while `errors.Is` asks the error itself.
