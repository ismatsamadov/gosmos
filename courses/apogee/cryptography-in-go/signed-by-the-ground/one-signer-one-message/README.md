---
slug: one-signer-one-message
title: One Signer, One Message
type: code_tests
position: 6
difficulty: 7
topics: crypto, goroutines, methods
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Friday morning, Ingrid's night log has a line for you: "02:14. Second command to Pipit-2 refused. Signature mismatch." The program that signs and sends commands, uplinkd, got a new `Signer` for its test on Pipit-2. It keeps one HMAC for speed, and it signed the night's second command together with the first.

A `hash.Hash` keeps everything written to it, and `Sum` doesn't clear it: what you write next is added to what came before. Start a new HMAC for every message, and each signature covers its own command only. No two calls share a hash then, so many goroutines can sign at once:

```go
h := sha256.New()
io.WriteString(h, "PING")
first := h.Sum(nil)
io.WriteString(h, "PING")
fmt.Println(bytes.Equal(first, h.Sum(nil)))
// false
```

## Your task

Make `Sign` start its own `hmac.New(sha256.New, s.key)` for every command, and remove the `mac` field and the part of `NewSigner` that fills it. The test signs two commands, then 400 at once from 8 goroutines.
