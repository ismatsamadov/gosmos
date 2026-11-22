Read the manifest with a `bufio.Scanner`, counting lines from 1, and split each line with `strings.Fields`. Three fields, a digest that `parseDigest` accepts and a size that `strconv.ParseInt` reads make a good line; anything else ends the check with the error.

For each file, `fsys.Open` fails with an error that `errors.Is` matches to `fs.ErrNotExist` when it's missing. Otherwise compare its size, from `fs.Stat` or the open file's `Stat`, before you read it, and only then hash it with `fingerprintAs` and the digest's own hash.

`fingerprintAs` returns the same `<hash>:<hex>` text the manifest holds, so compare the two strings: digests of public files need no constant-time comparison. Close every file you open.
