---
slug: extract-once-expand-many
title: Extract Once, Expand Many
type: code_tests
position: 5
difficulty: 8
topics: crypto, maps
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Wednesday. Thursday's passes carry every satellite's two keys up at once, so uplinkd derives fourteen keys from one master in one go.

`hkdf.Key` works in two steps every time you call it: `hkdf.Extract` turns the secret and the `salt` into one pseudorandom key, then `hkdf.Expand` turns that and the `info` into the key you asked for. When one master gives many keys at once, you can extract once and expand each, and the keys come out the same:

```go
prk, _ := hkdf.Extract(sha256.New, s, salt)
a, _ := hkdf.Expand(sha256.New, prk, "logs", 16)
b, _ := hkdf.Key(sha256.New, s, salt, "logs", 16)
fmt.Println(bytes.Equal(a, b)) // true
```

## Your task

Make `fleetKeys` extract once from the master's secret and `salt`, then expand a 32-byte key for `sign|<sat>` and one for `seal|<sat>` for every satellite in `satIDs`, in a map keyed by that info. The test reads `main.go`: one call to `hkdf.Extract`, outside any loop, and none to `hkdf.Key`.

> [!NOTE]
> Go's documentation suggests `hkdf.Key` for most code, several keys included, and `Extract` only when one extracted key must feed many `Expand` calls, as in this task.
