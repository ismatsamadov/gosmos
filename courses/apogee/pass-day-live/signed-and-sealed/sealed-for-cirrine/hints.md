`aes.NewCipher` takes the 32-byte seal key and `cipher.NewGCM` wraps it; its `Seal` method, given `nil`, the nonce, the plaintext and the associated data, returns the sealed bytes with the 16-byte tag at their end. Fill a new 12-byte nonce for every seal, a resent line aside, with `rand.Read` from `crypto/rand`.

`SEAL 2 <sat> <seq> <key ID>` starts the line, and the same five values joined by bars are the associated data. The key ID is the master the seal key came from, and the seal key is derived like the sign key, with `seal|<sat>` as the info.

`command failed authentication` from a Cirrine satellite means the associated data or the key is off: Cirrine-1's first PING has `SEAL|2|cirrine-1|9713|k-2028-10`, so compare yours field by field, and check the key's info. If a nonce sealed two lines, each seal must draw its own: a generator seeded with a number or the clock repeats itself.
