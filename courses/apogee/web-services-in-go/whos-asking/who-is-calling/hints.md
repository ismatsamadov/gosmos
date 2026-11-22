`r.Header.Get("Authorization")` gives the whole header, or `""` when there is none. `strings.Cut` splits it at its first space into the scheme and the rest.

Compare the scheme you found with `scheme` using `strings.EqualFold`, which ignores letter case. The rest is the credential, and an empty one counts as none.

If `key furrowlight...` is refused, you are comparing the scheme with `==` or `strings.HasPrefix`. If `Key` alone gets through, check that the credential isn't empty before you return it; every refusal returns `""` and `unauthorized("missing credentials")`.
