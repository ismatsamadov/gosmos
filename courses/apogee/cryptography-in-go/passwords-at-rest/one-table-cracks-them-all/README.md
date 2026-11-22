---
slug: one-table-cracks-them-all
title: One Table Cracks Them All
type: code_tests
position: 0
difficulty: 7
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Monday 31 July. On Sunday the Fishbowl's spare laptop went missing from the night-shift desk, with a copy of the ops console's user file: 26 accounts, each password kept as its SHA-256 digest. Mirela: "Assume whoever has it has a fast computer and a list of common passwords."

The same password always gives the same digest, so a table of common passwords and their digests, made once, reads at a glance every password record that holds one of them:

```go
table := map[[32]byte]string{} // digest → pw
for _, pw := range []string{"letmein", "orbit7"} {
	table[sha256.Sum256([]byte(pw))] = pw
}
stolen := sha256.Sum256([]byte("orbit7"))
fmt.Println(table[stolen]) // orbit7
```

A salt spoils the table: hash 16 random bytes and the password together, and keep the salt beside the result. Equal passwords then make different records, which no table made in advance can match.

## Your task

Make `hashPassword` take a fresh salt from `newSalt` and return this record, where the digest is SHA-256 of the salt's 16 bytes followed by the password's:

```text
salted$<32 hex salt>$<64 hex digest>
```

The test plays whoever has the laptop.
