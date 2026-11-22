For `newKey`, make a 32-byte slice, fill it with `rand.Read` as you did the nonce, and put it in the `Secret` field of the key you return.

`hex.DecodeString` refuses only what isn't hex. Check the length of `text` before decoding it: anything but 64 characters is an error, which `fmt.Errorf` can make.

If the 62-digit line or the empty one still gets through, the length check comes after the decode or compares with the wrong number: 32 bytes are 64 hex digits.
