---
slug: mirelas-random-kit
title: Mirela's Random Kit
type: code_tests
position: 7
difficulty: 7
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Friday 14 July. At 16:00 Mirela makes `k-2028-07` with your helpers, and before then she wants them in one file that every later tool takes its secrets from: "One file to review. One place secrets come from."

The kit makes four kinds of secret, each from `crypto/rand`. One is new: 16 random bytes kept beside each password, called a salt, for the ops console's new password records. The pattern is the one you know, for any size:

```go
// newAPIKey returns 16 random bytes in hex,
// as in a customer's API key.
func newAPIKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
```

## Your task

Write the four helpers in `random.go`:

```text
newKey(id)         Key{ID: id}, 32 random bytes
newNonce()         12 random bytes
newSalt()          16 random bytes
newSessionToken()  rand.Text()
```

> [!NOTE]
> The test reads every `.go` file of the lesson: `random.go` must import `crypto/rand`, and no file may import `math/rand` or `math/rand/v2`.
