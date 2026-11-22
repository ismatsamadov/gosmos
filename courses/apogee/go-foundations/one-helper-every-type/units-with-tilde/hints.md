Read the compiler's message to the end: it says what `Real` is probably missing.

`~float64` admits every type whose underlying type is `float64`. `KmPerHour` is built on `float32`, so `Real` needs a second term for it, joined with `|`.

`Real` keeps its shape, `type Real interface { ... }`. Only the line inside changes: two terms, each starting with `~`, with `|` between them.
