`mlkem.NewEncapsulationKey768(ek)` returns the receiver's key and an error; return the error at once, with nil for the key and the ciphertext.

`Encapsulate` takes no arguments and returns the shared key, then the ciphertext, in the order `wrapFor` returns them. In `unwrap`, `dk.Decapsulate(ciphertext)` returns the key and an error, ready to return.

If `wrapFor` reports a 1088-byte key and a 32-byte ciphertext, the two results are swapped: `Encapsulate` returns the key first and the ciphertext second.
