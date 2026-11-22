---
slug: open-in-the-right-order
title: Open in the Right Order
type: code_tests
position: 6
difficulty: 8
topics: crypto, errors
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T19:06:42Z
---
Thursday 24 August. On board Cirrine-1, the program that runs the satellite, its flight software, checks every sealed command in Mirela's order. Its stand-in on the ground, which every test talks to, must do exactly the same before tomorrow's pass.

Checks that read only the header come first and `Open` comes last. A number already seen is refused before any work is spent on it, and a number is recorded only after `Open` succeeds, so a forgery never moves the replay guard. The draft records the number before opening: one forged line numbered 18,446,744,073,709,551,615 makes every real command afterwards look like a replay. Each refusal wraps its error, for `errors.Is`:

```go
err := fmt.Errorf("%w: %d, last %d",
	ErrReplay, 4471, 4471)
fmt.Println(err)
// replayed sequence number: 4471, last 4471
fmt.Println(errors.Is(err, ErrReplay)) // true
```

## Your task

Put `Accept`'s checks in this order, and set `s.Last` only once `openCommand` has succeeded:

```text
1 parseSealed          ErrMalformed, ErrVersion
2 satellite is s.ID    ErrWrongSat
3 s.keyFor(key ID)     ErrUnknownKey, ErrRevoked,
                       ErrRetired
4 number above s.Last  ErrReplay
5 openCommand          ErrForged
```
