---
slug: quick-check-keys
title: Quick Check: Keys
type: quiz
position: 3
difficulty: 8
topics: http, security
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
On Friday Yusra sends the beta customers their API keys, and the service will keep only their hashes. Before she does, Mirela has three questions about this check and about the ways the service refuses a caller:

```go
sum := sha256.Sum256([]byte(key))
got := hex.EncodeToString(sum[:])
ok := subtle.ConstantTimeCompare(
	[]byte(got), []byte(c.KeyHash)) == 1
```

## Your task

Answer Mirela's three questions: why the comparison isn't `==`, what a stolen copy of the customer table is worth, and what Brinewatch should get for trying to book a pass of Cirrine's satellite.
