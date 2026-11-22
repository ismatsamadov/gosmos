Wrap `r.Body` in an `io.LimitReader` of 64 KiB, `64<<10` bytes, before `io.ReadAll`: a longer body then arrives cut short and can't match its signature. `strconv.ParseInt` returns an error for a missing header, which reads as an empty string.

`now.Unix() - ts` is negative for a time ahead of yours. Refuse it below -300 or above 300, so that exactly 300 either way is still accepted.

`hmac.Equal` takes two byte slices: the header's value and the hex string `signature` returns. Decode with `json.Unmarshal` only after it says true, and return `Event{}` with every error, a failed decode included.
