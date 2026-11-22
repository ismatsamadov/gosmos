`strings.Split(tok, ".")` gives the parts; anything but three is `invalid token`. The first two parts joined by a dot are the claims the signature was made over.

Decode the third part with `hex.DecodeString`, and treat an error as `invalid token`. Compute `tokenMAC` over the claims and compare the two with `hmac.Equal`.

Only after the signature checks out, parse the expiry with `strconv.ParseInt` and turn it into a time with `time.Unix(seconds, 0)`. The token is valid only while `now.Before` that time. If an edited token comes back `token expired`, you are checking the expiry before the signature.
