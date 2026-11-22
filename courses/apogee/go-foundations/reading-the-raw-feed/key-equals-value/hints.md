`strings.Cut` hands back three results: before, after and found. Which one do you compare with `key`, and which one do you return?

Call `strings.Cut(f, "=")` inside the loop; the part before the `=` is the field's key, and the part after it is everything else, even a second `=`.

A field like `thistle` has no `=`: `Cut` returns it whole as the first result, with found `false`. Check found as well, or `field(line, "thistle")` reports a key that isn't there.
