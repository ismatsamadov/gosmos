---
slug: who-is-calling
title: Who Is Calling?
type: code_tests
position: 0
difficulty: 8
topics: errors, http, security, strings
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Mirela's review of the booking service opens with one line: "Today anyone who can reach this service can read every booking. Who is calling?"

A request names its caller in its `Authorization` header: a scheme saying what kind of credential follows, one space, then the credential. Customers' programs will send the secret Yusra gives each of them, an API key, as `Key <api key>`, and later tokens as `Bearer <token>`. Schemes match in any letter case:

```go
cmd, rest, _ := strings.Cut("stow CWR", " ")
fmt.Println(cmd, "|", rest) // stow | CWR
fmt.Println(strings.EqualFold("STOW", "stow"))
// true
```

## Your task

`credentials` returns `""` and nil. Return the credential that follows `scheme` in `r`'s `Authorization` header. When the header is missing, names another scheme or has nothing after its space, return `""` and `unauthorized("missing credentials")`, the 401 that says "we don't know who you are".

> [!WARNING]
> `strings.HasPrefix(h, "Key ")` turns away `key furrowlight...`, a header that is fine.
