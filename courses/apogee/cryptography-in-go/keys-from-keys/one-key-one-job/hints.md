The info is now built from both arguments: the purpose, a bar, then the satellite's ID, so Cirrine-1's two keys come from `sign|cirrine-1` and `seal|cirrine-1`.

Check both arguments before deriving anything. A purpose is acceptable only if it is `sign` or `seal`, and `slices.Contains`, once `"slices"` is imported, tells you whether `fleet` holds the satellite. `fmt.Errorf` can name the value that was refused.

If `wrap` or `thistle` still gets a key, the checks come after the call to `hkdf.Key` or don't return. If the `seal` keys fail, the info still says `sign` whatever the purpose.
