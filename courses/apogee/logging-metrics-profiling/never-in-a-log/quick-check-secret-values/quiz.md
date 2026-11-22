## Without `blank`, what does the `command sealed` line write for `k`?

- [ ] `"key":"k-2028-10 (secret hidden)"`, from its `String` method
  > The JSON handler never calls `String`; fmt and the text handler do.
- [x] Its `ID` and its `Secret`, the secret in base64
  > With no `LogValue` to call, the JSON handler writes `k` as `encoding/json` would: every exported field, a `[]byte` in base64.
- [ ] Nothing for `k`: the JSON handler skips structs
  > It writes any value; a struct becomes an object of its exported fields.

## Which attributes of the `command sealed` line is `blank` called for?

- [ ] Only `key`, the one it changes
  > It is called for every attribute; `blank` itself decides which ones to change.
- [ ] `seq` and `key`, never the built-in ones
  > The built-in ones pass through it too: that is how `dropTime` has left out the time since chapter 1.
- [x] `time`, `level`, `msg`, `seq` and `key`, each with the groups it sits in
  > A handler calls `ReplaceAttr` for every attribute that isn't a group, the built-in ones included, with the names of the groups around it.

## The text of `err` holds Furrowlight's old API key. What does `blank` do about it?

- [x] Nothing: the secret sits inside the value under `err`
  > `blank` looks only at an attribute's key, never inside its value, so a secret in an error's text passes it.
- [ ] Replaces the text, since it holds a secret
  > `blank` looks at each attribute's key, and this one's key is `err`.
- [ ] Drops the whole record, since it holds a secret
  > `ReplaceAttr` rewrites one attribute at a time; it can't drop a record.
