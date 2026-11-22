`strings.Cut` already tells you whether the text had a `:`: use its third result, and return the first error at once when it is false.

Compare the name with the two allowed names exactly; `fmt.Errorf` with `%q` writes the name back in quotes. For the digits, check that there are 64 of them and that none is outside `0` to `9` and `a` to `f` before decoding.

`hex.DecodeString` accepts uppercase letters, so it can't be the only check. If the trailing space slips through, your check of the digits' count comes after the decode, or trims the text first.
