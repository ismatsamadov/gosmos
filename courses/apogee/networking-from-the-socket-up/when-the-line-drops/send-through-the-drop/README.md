---
slug: send-through-the-drop
title: Send Through the Drop
type: code_tests
position: 7
difficulty: 9
topics: context, errors, networking
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Thursday. Before tomorrow's drill, Sunniva cuts the test line to the simulator the moment a STOW, which turns a dish straight up and holds it, leaves dishd: "Calm. Cutting at your STOW. Stow once, please."

`Send` brings the chapter together. The given `connect` redials, says HELLO, learns the count and calls `resume`, so a nil `l.inflight` afterwards means the command in flight has run. The given `settle` judges an answer: ACK or NAK `duplicate` is nil, any other NAK its `*NakError`.

```go
ans, err := l.client.Do(ctx, stow)
fmt.Println(ans.Type, err) // 0 EOF
fmt.Println(lost(err))     // true
```

## Your task

`Send` in `link.go` numbers the command and sends it once. Keep it in `l.inflight` and, on a lost line, close `l.client` and `connect`, then return nil if it ran or send the same frame again; let `settle` judge the answer.

> [!WARNING]
> Number it once, before the first try: under a new number, a command that already ran runs again.
