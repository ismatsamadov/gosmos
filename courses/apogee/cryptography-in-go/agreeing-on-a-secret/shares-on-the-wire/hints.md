Decode first, then refuse the text if decoding failed, if the result isn't 32 bytes, or if `hex.EncodeToString` of the result doesn't give back exactly the text that came in. Capitals fail that last check, since `EncodeToString` always writes lowercase.

In `agreeHex`, the line that calls `mine.ECDH` drops its error with `_`. Return both of its results instead.

If the test of 64 zeros still finds no error, `agreeHex` is still dropping the error from `ECDH`; `parseShare` is right to accept the zeros, since they are 64 well-formed digits.
