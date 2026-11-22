---
slug: keys-never-print
title: Keys Never Print
type: code_tests
position: 4
difficulty: 8
topics: formatting, interfaces, methods
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Wednesday. Last night's uplinkd log holds all 32 bytes of `k-2028-07`: a debug line printed a key with `%v`. The log never left that machine, so Mirela deletes it and opens a review: "Make the keys refuse to print."

`fmt` calls a value's `String` method for `%v`, `%+v` and `%s`, and its `GoString` method for `%#v`, so methods that leave the secret out keep it out of what those verbs print. But `fmt` can't call methods on a value it reaches through an unexported field, and prints that value's fields instead:

```go
type Temp float64

func (t Temp) String() string { return "warm" }

fmt.Println(Temp(41.5))             // warm
fmt.Println(struct{ T Temp }{41.5}) // {warm}
fmt.Println(struct{ t Temp }{41.5}) // {41.5}
```

## Your task

In `main.go`, give `Key` a `String` and a `GoString` method. The ring keeps its versions in an unexported field, so give `*KeyRing` its own `String`, and a `GoString` that returns the same text. For the test's keys, in the ring's order:

```text
k-2028-07 (secret hidden)
Key{ID: "k-2028-07", Secret: hidden}
key ring: k-2028-07, k-2028-10
```

> [!WARNING]
> `%d` skips `String` and prints a key's bytes as numbers, so log keys with `%v`.
