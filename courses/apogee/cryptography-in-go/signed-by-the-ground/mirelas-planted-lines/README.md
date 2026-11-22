---
slug: mirelas-planted-lines
title: Mirela's Planted Lines
type: code_tests
position: 7
difficulty: 8
topics: crypto, errors, strings
created: 2026-11-07T07:46:48Z
modified: 2026-11-22T10:06:32Z
---
Friday 28 July, 16:00. Mirela has written twelve command lines by hand, two good ones and ten bad in ten different ways, and she feeds them to the uplink's checker: "Accept two. Refuse ten. And show me the guard never moved."

A checker that fails closed treats every doubt as a refusal, and no refusal changes anything. The line is the one Pipit-2 has always read: `CMD`, the command's fields, then one `sig=` field, last. `strings.Fields` splits it:

```go
f := strings.Fields("CMD pipit-1 7 PING sig=ab12")
sig, ok := strings.CutPrefix(f[len(f)-1], "sig=")
fmt.Println(len(f), sig, ok) // 5 ab12 true
```

## Your task

`checkLine` refuses every line. Make it parse the line into a `Command`, refusing anything but `CMD`, a satellite, a sequence number, an operation, any arguments and a last `sig=` field of 64 lowercase hex digits. Then hand the command to `g.Accept`, which verifies before it records, and return its error as it is.
