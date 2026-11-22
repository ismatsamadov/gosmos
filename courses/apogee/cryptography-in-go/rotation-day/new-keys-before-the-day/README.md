---
slug: new-keys-before-the-day
title: New Keys Before the Day
type: code_tests
position: 3
difficulty: 10
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T19:06:42Z
---
Every satellite must hold `k-2028-10` before Rotation Day, loaded on its next pass, starting with Pipit-1 tonight. Mirela's condition: "The new keys never cross the air readable."

A satellite's new keys travel inside a LOAD_KEY command, sealed under a key the satellite already holds. Its plaintext, 148 bytes for `k-2028-10`, is `LOAD_KEY|<id>|<sign hex>|<seal hex>`, and the satellite installs the keys only once the command opens. `sealCommand` seals any command:

```go
c := Command{SatID: "pipit-1", Seq: 812,
	Op: "SET_MODE", Args: []string{"SAFE"}}
line, err := sealCommand(ring, "k-2028-07", c)
fmt.Println(line[:28], err)
// SEAL 2 pipit-1 812 k-2028-07 <nil>
```

## Your task

In `loadKeyCommand`:
1. Derive `satID`'s `sign` and `seal` keys from `next.Key` with `satKey`.
2. Build the `LOAD_KEY` command numbered `seq`, its arguments `next.ID` and the two keys in hex.
3. Seal it under the version `signingKey` picks one second before `next.Activate`.
