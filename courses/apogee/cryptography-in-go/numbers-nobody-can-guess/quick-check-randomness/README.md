---
slug: quick-check-randomness
title: Quick Check: Randomness
type: quiz
position: 2
difficulty: 6
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Wednesday. Mirela rehearses Friday with a practice key from your helper: printed, folded into an envelope and locked in the safe. Before the safe closes, she writes on the whiteboard the kind of key maker she never wants to see:

```go
r := rand.New(rand.NewPCG(7, 7)) // math/rand/v2
key := make([]byte, 32)
for i := range key {
	key[i] = byte(r.Uint32())
}
```

## Your task

Answer Mirela's three questions: what this code makes when it runs again tomorrow, what `rand.Read` from `crypto/rand` promises, and how many different values a 32-byte key can take.
