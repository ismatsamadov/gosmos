A fresh key depends on nothing the `Backhaul` remembers, so the struct needs no field at all: delete `static` and the lines that read or set it.

`ecdh.X25519().GenerateKey(rand.Reader)` returns a new key and an error, which `NewSession` can return as they come.

If `TestSharesDiffer` still sees one share twice, some path still returns a key made earlier: every call must reach `GenerateKey`. If `TestKeepsNothing` names a field, the struct still keeps it, even when every call makes a new key: a thief would copy the last one.
