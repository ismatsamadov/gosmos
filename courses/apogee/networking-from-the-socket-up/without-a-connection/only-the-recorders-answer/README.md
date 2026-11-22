---
slug: only-the-recorders-answer
title: Only the Recorder's Answer
type: code_tests
position: 5
difficulty: 9
topics: networking, strings
created: 2026-11-13T06:15:46Z
modified: 2026-11-13T16:03:24Z
---
Thursday 15 June. The recorder now answers questions on the port it sends from, its live port: `STATUS 7` gets `STATUS 7 RECORDING furrow-a/YST/18781`, the asker's number first. Mirela: "Anyone can send the live view anything. Show me the answer you refuse."

A socket with no connection takes datagrams from anyone. A datagram answers question 7 only when it comes from the recorder's address and starts with `STATUS 7` and a space; a late answer to question 6, or one to question 70, belongs to another question. Compare addresses by `String()`, as two values of one address are often different pointers.

```go
ip := net.IPv4(127, 0, 0, 1)
a := &net.UDPAddr{IP: ip, Port: 7}
b := &net.UDPAddr{IP: ip, Port: 7}
fmt.Println(a == b, a.String() == b.String())
// false true
```

## Your task

`answerTo` in `answer.go` accepts nothing. When `from` is `recorder` and `d` starts with `STATUS <n>` and a space, with some text after, return that text and true; the given `question(7)` returns `STATUS 7`.

> [!NOTE]
> An address can be forged: this check keeps out strays and mistakes, not someone set on lying.
