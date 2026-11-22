## What is `n`?

- [ ] 2
  > 2 is the number of characters, but `len` counts bytes, and `°` takes two of them in UTF-8.
- [x] 3
  > `°` takes two bytes in UTF-8 and `C` takes one, and `len` counts bytes.
- [ ] 4
  > A string literal holds UTF-8, not two bytes per character: `C` takes a single byte.

## What is `b`?

- [ ] The rune `°`
  > Indexing a string gives a byte, not a rune; `for range` is what hands you runes.
- [ ] The string `"°"`
  > Indexing gives one byte, not a one-character string.
- [x] A byte: the first of the two that encode `°`
  > `s[0]` is the string's first byte, of type `byte`, and here it is only half of `°`.

## What happens if the code adds the line `s[0] = 'x'`?

- [x] It doesn't compile, because strings can't be changed
  > Go strings are immutable, so the compiler refuses any assignment to `s[0]`.
- [ ] `s` becomes `"xC"`
  > Strings can't be changed in place, and `s[0]` is only one byte of `°`, not the whole character.
- [ ] It compiles, then panics when it runs
  > The compiler catches it, so the program never gets to run.
