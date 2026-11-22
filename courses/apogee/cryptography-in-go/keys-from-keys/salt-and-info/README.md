---
slug: salt-and-info
title: Salt and Info
type: code
position: 1
difficulty: 7
topics: crypto, formatting
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Tomás wants the fleet's new keys side by side before any of them goes up: "Line them up like a frequency plan, four bytes each. If two match, I want to see it on the bench, not in orbit."

HKDF mixes its `salt` and `info` into every byte it returns. The `salt` is fixed and public, `apogee uplink keys v1` for all of Apogee's uplink keys, and keeps them apart from any other use of the same master. The `info` names one key, so each name gives a key unrelated to all the others. Change either one, and the key changes:

```go
m := []byte("demo master, never a real key 32")
s1, s2 := []byte("v1"), []byte("v2")
a, _ := hkdf.Key(sha256.New, m, s1, "x", 4)
b, _ := hkdf.Key(sha256.New, m, s2, "x", 4)
fmt.Printf("%x %x\n", a, b) // 0b01b118 3acd3d41
```

## Your task

`main` derives every line's key with the same info, so all four lines show the same bytes. Pass each line's label as the info, so that it prints:

```text
sign|pipit-2   82f94563
sign|furrow-a  00a1988a
sign|cirrine-1 18668bfe
seal|cirrine-1 427a8ae7
```
