---
slug: a-homemade-signature
title: A Homemade Signature
type: code_tests
position: 0
difficulty: 7
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Monday 24 July. The program on Pipit-3 itself, its flight software, accepts a command when the short proof sent with it, its tag, equals `sha256(key + command)`. Mirela makes a valid tag for a command with `|REBOOT` on the end without ever seeing the key: "A hash of a secret is not a signature."

A SHA-256 digest is the whole state the hash ended in, so anyone holding one can carry on hashing from there: the key stays inside, and the message grows. That trick is a length extension. An HMAC hashes the key with the message, then hashes that digest with the key again, which leaves nothing useful to carry on from:

```go
mac := hmac.New(sha256.New, []byte("demo key"))
io.WriteString(mac, "cirrine-1|7|PING")
fmt.Printf("%x\n", mac.Sum(nil)[:4]) // d6ae0ed9
```

## Your task

`tag` returns `sha256.Sum256` of the key followed by the command. Make it return the command's HMAC-SHA256 under the key. The test plays Mirela, extending the tag of `pipit-3|212|SET_MODE|SAFE`.
