---
slug: make-every-guess-cost
title: Make Every Guess Cost
type: code_tests
position: 1
difficulty: 7
topics: constants, crypto
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
A salt stops the table, not the guessing. Mirela's own laptop tries a billion salted SHA-256 guesses a second, so the whole list of 1,000 common passwords costs it a microsecond per record.

The cure is to make every guess expensive. `pbkdf2.Key` from `crypto/pbkdf2`, a standard recipe called PBKDF2, mixes the password and the salt through a hash over and over, as many rounds as its iteration count says, and returns a key of the length you choose. At 600,000 iterations the console spends about 60 ms checking a login, and anyone guessing pays those 600,000 rounds for every single guess:

```go
s := []byte("example salt 16b")
k, err := pbkdf2.Key(sha256.New, "letmein", s,
	4096, 16)
fmt.Println(len(k), err) // 16 <nil>
```

## Your task

1. Set the constant `iterations` to 600000.
2. Make `derive` return the 32-byte key from `pbkdf2.Key`, with `sha256.New`, the password, the salt and `iterations`.

> [!NOTE]
> One round costs about as much as two SHA-256 guesses, so at 600,000 iterations Mirela's billion guesses a second become fewer than a thousand.
