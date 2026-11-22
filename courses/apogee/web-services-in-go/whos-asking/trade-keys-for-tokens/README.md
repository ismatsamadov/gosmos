---
slug: trade-keys-for-tokens
title: Trade Keys for Tokens
type: code_tests
position: 6
difficulty: 8
topics: http, json, security, time
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Yusra's TL;DR from Hollis: "Furrowlight's program will ask for a new token every hour, from Monday." It sends its API key to `POST /tokens` and carries the token on every other call.

The handler checks the key, then answers with a `TokenResponse`: the token and its `expires_at`. A refusal also names the scheme to authenticate with, in a `WWW-Authenticate` header: here, `Key`. A token keeps whole seconds only, so the expiry is truncated to the second before it is signed and sent:

```go
t := time.Date(2028, 1, 5, 9, 30, 15,
	600_000_000, time.UTC)
fmt.Println(t.Truncate(time.Second))
// 2028-01-05 09:30:15 +0000 UTC
```

## Your task

`issueToken` answers 501. Read the `Key` credential with `credentials` and check it with `checkKey`, refusing a key it doesn't know with `unauthorized("invalid api key")`. Answer either failure through `writeError`, with `WWW-Authenticate: Key`. Otherwise answer 200 with a `TokenResponse` for a token that expires an hour after `s.Now()`, truncated to the second.
