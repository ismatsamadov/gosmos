---
slug: upgrade-at-login
title: Upgrade at Login
type: code_tests
position: 5
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Ingrid reset every password on Monday, and the console stored the new ones the only way it knew: old `sha256$` records. Twelve engineers are on leave until after Friday, so their records stay old until they log in again.

The console never keeps a password, so a login is the one moment it holds the password again. When the password matches a weak record, an old `sha256$` one or one under 600,000 iterations, the console makes a new record from the password it now holds, and saves it in place of the old one. `strings.CutPrefix` tells the two kinds apart:

```go
rest, ok := strings.CutPrefix("CMD 118", "CMD ")
fmt.Println(ok, rest) // true 118
```

## Your task

1. In `checkStored`, check an old `sha256$` record yourself: compare the hex SHA-256 of the password with the record's 64 digits using `subtle.ConstantTimeCompare`. Report every old record, and every record under 600,000 iterations, as weak.
2. In `login`, when the password matched a weak record, return a new one from `newRecord`.

> [!NOTE]
> The tests use the kiosk's record from the stolen file. Its password, `fishbowl`, is no secret anymore.
