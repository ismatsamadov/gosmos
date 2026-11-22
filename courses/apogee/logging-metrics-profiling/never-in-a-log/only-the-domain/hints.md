`strings.Cut(s, "@")` gives the part before the first `@`, the rest after it, and whether there was an `@` at all.

Return `[redacted]` when there is no `@`, when either part is empty, or when the rest still holds an `@`. Otherwise join the first character of the part before `@`, then `***@`, then the domain.

If `hollis@vey@furrowlight.example` comes back masked, look for a second `@` in the part after the first. Take the first character as a rune rather than a byte, so a letter of more than one byte stays whole.
