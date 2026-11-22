Looping over `name` by byte index splits `ñ` into two pieces. Which kind of loop hands you whole runes?

Ranging over `name` gives you one rune per pass, and a `switch` with no value can sort each rune by comparing it with rune literals such as `'a'` and `'z'`.

Four kinds of rune, four cases: a to z (write its upper case; `unicode.ToUpper` makes it), A to Z with digits and `-` (write the rune as it is), a space (write a hyphen), and `default` for the rest (write `_`).
