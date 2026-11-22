---
slug: a-list-of-the-revoked
title: A List of the Revoked
type: code_tests
position: 4
difficulty: 10
topics: crypto, errors, time
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Tuesday 17 October, 10:12 UTC. The Cape Wren contractor calls the Fishbowl: the old station computer she took away for recycling still holds its disk, and on it the key of certificate 1005. Cape Wren switched to 1105 at 01:00, but 1005 stays valid until December.

A revocation list names the certificates a certificate authority withdraws before they expire, numbered, dated and signed by that authority. A checker refuses a list whose signature fails or whose dates don't cover the moment, rather than trust it:

```go
list, err := x509.ParseRevocationList(der)
if err != nil {
	return err
}
fmt.Println(list.Number, list.ThisUpdate)
// 1 2028-10-17 10:12:00 +0000 UTC
```

## Your task

1. Make `revoke` sign list `number` of `entries` with `ca`: `ThisUpdate` is `now`, `NextUpdate` 7 days later.
2. Make `isRevoked` parse the list, check it with `CheckSignatureFrom(ca)`, fail unless `now` is from `ThisUpdate` to `NextUpdate`, both included, then look for `cert`'s serial number.
