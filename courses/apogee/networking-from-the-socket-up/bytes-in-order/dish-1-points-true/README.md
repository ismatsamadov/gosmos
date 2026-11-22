---
slug: dish-1-points-true
title: Dish 1 Points True
type: code_tests
position: 8
difficulty: 8
topics: errors, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Friday 31 March, 06:50. Cirrine-1's 07:02 pass will be the first customer pass on Dish 1's new controller, and Ingrid's go/no-go poll needs the dish pointed. Every command has to say how it ended.

A command is answered by an ACK, the controller's "done", or by a NAK, its refusal, with a code and a reason. Both carry the command's sequence number, so an ACK with another number answers some other command. `errors.AsType` finds an error of one type anywhere in a chain:

```go
err := c.command(unstow)
if ne, ok := errors.AsType[*NakError](err); ok {
	fmt.Println(ne.Code, ne.Reason)
	// 3 wind lock
}
```

## Your task

`command` sends the frame and returns. Read the answer with `recv`: an ACK of the same seq is `nil`, a NAK is the given `nakError(ans)`, and anything else is an error naming its type and both seqs:

```text
answer to seq 119: type 0x90, seq 118
```

Return `recv`'s own errors as they are.
