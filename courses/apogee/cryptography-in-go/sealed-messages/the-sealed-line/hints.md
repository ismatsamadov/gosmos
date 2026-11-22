`formatSealed` is one `fmt.Sprintf`: the word `SEAL`, then six verbs, `%d` for the version and the number, `%s` for the two IDs and `%x` for the nonce and the ciphertext, single spaces between them.

In `parseSealed`, split the line on `" "` and check the count first, since every later check reads a field by its index. On a line that starts with `SEAL`, a version `f[1]` other than `"2"` is an error wrapping `ErrVersion` with `%w`, followed by the version found. `strconv.ParseUint` with base 10 and 64 bits reads the number, and refuses `44a1`.

If a malformed case still gets through, list what it breaks: `f[0]` must be `SEAL`, `validID` must accept `f[2]` and `f[4]`, the number must be 1 or more, and `hexBytes` must accept the nonce at exactly 12 bytes and the ciphertext at 17 to 272. Any failure wraps `ErrMalformed`.
