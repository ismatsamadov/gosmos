Compare the attribute's whole key with the five names in one letter case: `strings.ToLower` of `a.Key` in a `switch`, or `strings.EqualFold` against each name. The groups don't change the answer, so `redact` can leave them alone.

For a key on the list, return a new attribute with the same key and the string `[redacted]`, made with `slog.String`; return every other attribute exactly as it came.

If `token_count` comes out blank, you matched part of a key: compare whole keys. If the key itself is missing from the line, you returned an empty `slog.Attr`, which drops the attribute instead of blanking it.
