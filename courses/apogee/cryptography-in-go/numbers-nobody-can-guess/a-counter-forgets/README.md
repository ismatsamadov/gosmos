---
slug: a-counter-forgets
title: A Counter Forgets
type: code_tests
position: 5
difficulty: 7
topics: crypto, methods
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Nkem sketches a simpler nonce maker for the program that sends Cirrine's commands up, uplinkd: count 0, 1, 2 and so on, so no two nonces are alike. Then she asks: "What does Cirrine-1 see after a restart?"

Random nonces need nothing remembered. A count kept in memory starts from zero whenever its program starts, so a second run hands out the first run's nonces again, under the same key:

```go
type Tally struct{ n int }

func (t *Tally) Next() int { t.n++; return t.n }

var mon, tue Tally // two runs
fmt.Println(mon.Next(), mon.Next()) // 1 2
fmt.Println(tue.Next())             // 1
```

## Your task

Make `Nonce` return 12 bytes from `rand.Read` instead of the count, and remove the `next` field. The test starts two `Noncer` values, as two runs of the program would, and compares their first nonces.

> [!NOTE]
> A nonce isn't secret: it goes up in the open beside its command. It must never repeat under one key; random 12-byte ones stay safe for 2^32^ commands per key.
