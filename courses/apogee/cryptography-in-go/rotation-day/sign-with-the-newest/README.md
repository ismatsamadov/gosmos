---
slug: sign-with-the-newest
title: Sign With the Newest
type: code_tests
position: 1
difficulty: 9
topics: crypto, errors, time
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
uplinkd must seal with `k-2028-10` from 00:00:00 exactly on Rotation Day, and never with a revoked key. Its draft signer takes the first version in the ring, whatever order the key file happened to list them in.

A signer uses the newest version, by `Activate`, that has activated and is neither retired nor revoked. With none, it refuses to sign rather than guess. Instants compare with `Before`, `After` and `Equal`:

```go
a := time.Date(2028, 10, 17, 0, 0, 0, 0, time.UTC)
b := a.Add(-time.Second)
fmt.Println(b.Before(a), a.After(b)) // true true
fmt.Println(a.Format(time.RFC3339))
// 2028-10-17T00:00:00Z
```

## Your task

Make `signingKey` return that version, whatever the ring's order, or else the error below, with `now` written in RFC 3339:

```text
no key to sign with at 2028-10-25T00:00:00Z
```
