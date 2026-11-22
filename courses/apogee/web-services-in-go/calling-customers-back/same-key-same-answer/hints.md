Hash the body with `sha256.Sum256`, then hold `s.mu` for the whole check and claim, so two requests with one key can't both find it free. The keys live in the map `s.keys`, and the time is `s.now()`.

A kept key still counts while `now.Sub(k.At) < 24*time.Hour`. Compare `k.Sum` with the hash first, then look at `k.Done`, and return the matching refusal with a nil `*kept`.

To claim the key, store a new `kept` holding the hash and the time under `id` before returning nil, nil: `finish` can only fill in the answer of a key that was claimed, so without it no retry is ever replayed.
