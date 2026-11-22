---
slug: one-key-opens-everything
title: One Key Opens Everything
type: code_tests
position: 0
difficulty: 7
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Monday 7 August, at the lab bench. Mirela reads `k-2027-03` out of the spare Pipit-2 computer kept there, signs a command for Cirrine-1 with it, and the command checks out: all seven satellites check their commands with that one key.

Each satellite needs keys of its own, made from one master key that only uplinkd holds. `hkdf.Key` from `crypto/hkdf` derives as many keys as you need from one random secret, by a standard recipe called HKDF: give it a hash, the secret, a `salt`, an `info` string naming the key and a length, and it returns that many bytes. A derived key gives away neither the master nor any other key:

```go
m := []byte("demo master, never a real key 32")
a, _ := hkdf.Key(sha256.New, m, nil, "logs", 4)
b, _ := hkdf.Key(sha256.New, m, nil, "backups", 4)
fmt.Printf("%x %x\n", a, b) // 7f7ab6b3 b9acf2c0
```

## Your task

Make `signKey` return 32 bytes from `hkdf.Key`, with `sha256.New`, the master's `Secret`, the given `salt` and the info `sign|` followed by the satellite's ID, and return an error for an empty ID. The test plays Mirela at the bench.
