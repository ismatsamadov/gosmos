---
slug: readable-by-brinewatch-alone
title: Readable by Brinewatch Alone
type: code_tests
position: 4
difficulty: 9
topics: crypto, security
created: 2026-11-13T03:46:39Z
modified: 2026-11-13T03:46:39Z
---
Friday, 17:30. Brinewatch's new receiving key passes `readRSA`: 3,072 bits. Two new secrets must now leave locked so that only that key's private half opens them: the one Brinewatch's programs prove themselves with, its API key, and the one that proves Apogee's messages are Apogee's, its webhook secret.

RSA can lock a short message to a public key. Its safe form, OAEP, mixes in random bytes, so the same secret never locks the same way twice, and takes a label: a ciphertext opens only under the label it was locked with, so the API key's can't pass for the webhook secret's. A 3,072-bit key's ciphertexts are 384 bytes, and it locks at most its size less twice SHA-256's 32 bytes, less 2:

```go
fmt.Println(pub.Size())            // 384
fmt.Println(pub.Size() - 2*32 - 2) // 318
```

## Your task

1. Make `encryptTo` lock `secret` with `rsa.EncryptOAEP`, `sha256.New()`, `rand.Reader` and `label` as bytes.
2. Make `decryptWith` open it with `rsa.DecryptOAEP`, the same hash and the same label.
