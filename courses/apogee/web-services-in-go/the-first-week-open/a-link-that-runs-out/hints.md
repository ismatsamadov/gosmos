`strconv.ParseInt(expires, 10, 64)` turns the expiry into the number `linkSig` takes; an error there is a bad link.

`hmac.Equal` compares two byte slices in constant time: give it `[]byte(sig)` and the bytes of the hex string `linkSig` returns for this path and expiry.

`time.Unix(exp, 0)` is the expiry as a time. A link still works while `now` is before it, so at that very second `now.Before` is false and the link has expired.
