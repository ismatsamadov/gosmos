`hkdf.Key` takes the hash constructor `sha256.New`, the secret, the salt, the info and the key's length, and returns the key and an error you can return as they come. Add `"crypto/hkdf"` to the imports.

The secret is `master.Secret`, the salt is the `salt` variable already in `main.go`, and the info is the text `sign|` followed by `satID`, so Pipit-2's is `sign|pipit-2`.

Check `satID` before deriving and return an error when it is empty. If `TestStolenKey` still fails, some path still returns `master.Secret` itself, which every satellite would share.
