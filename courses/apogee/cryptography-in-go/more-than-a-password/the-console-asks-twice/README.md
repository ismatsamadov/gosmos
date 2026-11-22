---
slug: the-console-asks-twice
title: The Console Asks Twice
type: code_tests
position: 8
difficulty: 10
topics: crypto, errors, security
created: 2026-11-13T06:18:24Z
modified: 2026-11-13T16:11:48Z
---
Friday 3 November, 15:30. At 16:00 the console stops taking a password alone, and Mirela wants the whole door in one place: `Gate`, with the reset beside it.

`gate.go` holds the old `Login`'s work in pieces: `admit` upgrades a weak record and starts a session, `setPassword` saves a new record, `signOut` ends an account's sessions. One rule is new. Six digits can be guessed, so after `maxMisses`, ten wrong codes in a row, an account takes no code at all until its owner resets the password. A code is tried only after the right password, so only someone holding the password can lock an account.

```go
misses := map[string]int{}
misses["mirela"]++
fmt.Println(misses["mirela"]) // 1
delete(misses, "mirela")
```

## Your task

1. `SignIn`: refuse a wrong password, through `g.console.Check`; then a locked account; then a code that neither `g.codes` nor `g.recovery` accepts, adding a miss. Every refusal is `errInvalidSignIn`. Otherwise clear the misses and `admit`.
2. `ResetPassword`: `Redeem` the token, then `setPassword`, `signOut` and clear the misses.
