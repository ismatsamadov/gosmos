Start from an empty `url.Values{}` and call its `Set` method once for each value that isn't empty, under the key the service expects: `sat`, `station` or `from`.

`Encode` returns an empty string when nothing was set. In that case return `/passes` alone; otherwise return `/passes?` followed by what `Encode` returned.

If only Hollis's case fails, compare its `from`: `Encode` writes the `+` as `%2B` and the `:` as `%3A`, and puts `from` before `sat` because it sorts the keys. A string glued by hand does neither.
