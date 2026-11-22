---
slug: the-console-signs-in
title: The Console Signs In
type: code_tests
position: 7
difficulty: 8
topics: maps, methods, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Friday 4 August. Every engineer has a new password, and at 16:00 the console's login opens again. Ingrid's checklist has one line left before then: your `Login`.

A login ties the week together. `Check` answers for the password: its `matches` now reads old `sha256$` records too, as your `checkStored` did, and spends one derivation on the dummy record for each, so an old record answers no faster than a new one. A weak record is replaced, in the user file through `c.save` and in `c.records`; then the browser gets a new session token. The console keeps only the token's SHA-256 in hex, so a stolen copy of its sessions holds nothing that logs anyone in. Every later request is checked the same way:

```go
sum := sha256.Sum256([]byte(token))
user, ok := sessions[hex.EncodeToString(sum[:])]
```

## Your task

Make `Login`:

1. Return `Check`'s error when it refuses.
2. When `needsUpgrade` says the user's record is weak, make a new one with `c.newRecord`, save it with `c.save` and put it in `c.records`.
3. Make a token with `newSessionToken`, map its hex SHA-256 to the user in `c.sessions`, and return the token.
