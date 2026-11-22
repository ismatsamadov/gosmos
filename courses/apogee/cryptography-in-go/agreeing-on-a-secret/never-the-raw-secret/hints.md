Add `"crypto/hkdf"` to the imports and remove the placeholder line. Check `len(shared)` first: anything but 32 is an error, which `fmt.Errorf` can make.

Build the info with `hex.EncodeToString` of each share, the dialer's first, each after a `|`. `hkdf.Key` then takes the hash's constructor, `sha256.New`, the secret, nil for the salt, the info and the length you want, 32, and returns the key and an error, both ready to return.

If the swapped shares still give the same key, the info leaves the shares out or puts them in a fixed order; if the vector differs, compare your info with the line in the readme character by character.
