The two lengths are `len` of the slices `main` already has. The sizes come from methods of `aead`: `aead.NonceSize()` and `aead.Overhead()`.

For the last line, copy `sealed` with `bytes.Clone`, change the copy's last byte (XOR it with 1), and open the copy with the same nonce and `nil` associated data. Keep the error `Open` returns and print it with `%v`; the plaintext it returns is empty.

If the last line prints `<nil>`, the copy you opened is still the original: change a byte of the clone, not of `sealed`, and pass the clone to `Open`.
