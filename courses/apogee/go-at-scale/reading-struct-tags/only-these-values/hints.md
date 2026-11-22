Add a `case "oneof", "len":` to the `switch`. For a `[]string`, loop over `v.Len()` and check each `v.Index(i)` the way you check one string: `checkRule` itself can do that.

A string's `Value` gives its text with `String()`. `strings.Fields` splits the allowed words, `slices.Contains` looks for the text among them, and `utf8.RuneCountInString` counts its characters.

If the accented code fails or a field of another type gets a quoted `<int Value>`, the case counts bytes with `len` or calls `String()` on a number: count runes, and return `misfit` unless the value is a string or a `[]string`.
