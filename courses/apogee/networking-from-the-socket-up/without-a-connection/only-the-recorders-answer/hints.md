Check the sender first: `from.String()` must equal `recorder.String()`. Comparing `from == recorder` compares two pointers, and they differ even for one address.

Build the prefix from `question(n)` and a space, so `STATUS 7 ` can't match the start of `STATUS 70`. `strings.CutPrefix` returns what follows a prefix and whether the prefix was there.

An answer needs some text after its prefix: return false when what follows is empty.
