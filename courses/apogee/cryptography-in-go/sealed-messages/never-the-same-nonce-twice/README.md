---
slug: never-the-same-nonce-twice
title: Never the Same Nonce Twice
type: code_tests
position: 2
difficulty: 8
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Wednesday. The draft seals every command with the same nonce, the fixed one the test values use. Mirela records two sealed commands, knows what the first one says, and reads the second without the key.

GCM hides a command by XORing it with a stream made from the key and the nonce. Seal two commands with the same pair, and XORing their sealed bytes cancels the stream, leaving the two commands XORed together; knowing one gives away the other, as far as the known one reaches:

```go
ks := byte(0x9c)             // one stream, twice
a, b := byte('S'), byte('D') // two first bytes
ca, cb := a^ks, b^ks         // what goes up
fmt.Printf("%c\n", ca^cb^a)  // D
```

## Your task

Make `sealCmd` take a new nonce from `newNonce()`, in `random.go`, for every command it seals. It returns the nonce with the sealed bytes, because the nonce goes up beside them, in the open. The test plays Mirela with two of your seals.

> [!WARNING]
> Under one key, a nonce must never repeat, not even once: one repeat was all Mirela needed.
