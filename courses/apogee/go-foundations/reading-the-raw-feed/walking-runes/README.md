---
slug: walking-runes
title: Walking Runes
type: code_tests
position: 4
difficulty: 5
topics: loops, runes, strings
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Paz's new radio log tags every entry with a call sign, the short name a radio station uses to identify itself. The radios accept only A to Z, digits, hyphens and underscores.

`for i, r := range s` walks a string one rune at a time: `r` is the rune and `i` the byte where it starts, so a two-byte rune still comes out once. A rune is a number, so you can compare it with rune literals in single quotes, like `r >= '0'`:

```go
for i, r := range "Tomás" {
	fmt.Printf("%d:%c ", i, r) // %c prints a rune
}
// 0:T 1:o 2:m 3:á 5:s
```

## Your task

Fill in `callSign`. For each rune of `name`, write to `b` with `b.WriteRune`: the letters a to z and A to Z in upper case, which `unicode.ToUpper` gives you, digits and hyphens as they are, a hyphen for a space, and `_` for anything else.

> [!NOTE]
> A string can't be changed in place: `name[0] = 'S'` doesn't compile. That's why `callSign` builds a new one in `b`, a `strings.Builder`.
