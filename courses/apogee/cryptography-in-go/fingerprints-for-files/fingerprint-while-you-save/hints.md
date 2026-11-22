Make the hash first with `sha256.New()`, then pass `io.MultiWriter(dst, h)` to `io.Copy` as its writer, so every piece goes to both.

After a copy without an error, `h.Sum(nil)` holds the digest of exactly the bytes that reached `dst`. Return `"sha256:"` joined to its hex, from `hex.EncodeToString`, and remove the placeholder line.

If the digest is right but the full folder's case fails, return `""` together with the copy's error rather than a digest of the part that was saved.
