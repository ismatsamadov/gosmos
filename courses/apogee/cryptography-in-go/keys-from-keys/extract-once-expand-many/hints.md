Call `hkdf.Extract` once, before any loop, with `sha256.New`, the master's secret and `salt`, after adding `"crypto/hkdf"` to the imports. What it returns is the pseudorandom key that every `hkdf.Expand` call starts from.

Make the map first, then loop over `satIDs` and, inside, over the two purposes `sign` and `seal`. Build each info as the purpose, a bar and the satellite, expand 32 bytes from it, and store the key in the map under that info.

If `TestExtractOnce` fails, its message says what it found: `hkdf.Key` anywhere, or `hkdf.Extract` more than once or inside a loop. Return the first error `Extract` or `Expand` gives.
