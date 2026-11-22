---
slug: the-pass-in-the-context
title: The Pass in the Context
type: code_tests
position: 0
difficulty: 7
topics: context, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
By Monday evening, you and Nkem have matched the four programs' lines about Cirrine-1's pass by hand, and by guesswork. The pass's name, its pass ID `cirrine-1/LAN/16402`, must travel with the work, so every line can carry it.

A `context.Context` goes along with a piece of work through every call it makes. `context.WithValue` returns a child context holding one value under a key, and `ctx.Value(key)` finds it there or in any context made from it. Make the key an unexported struct type: no other package can make one, so a value someone stored under the plain string `"pass_id"` is never taken for yours:

```go
type siteKey struct{}
ctx := context.WithValue(context.Background(),
	siteKey{}, "LAN")
site, ok := ctx.Value(siteKey{}).(string)
fmt.Println(site, ok) // LAN true
```

## Your task

Write `withPass`, which returns `ctx` with `id` stored under `passKey{}`, and `passFrom`, which gives the ID back with `true`, or `""` and `false` when the context holds no ID or an empty one.
