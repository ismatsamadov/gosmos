A `Set[T]` is a map, so its methods use it like one. How do you store a key in a map, and how do you ask whether a key is there?

A map of `struct{}` stores the empty struct value, written `struct{}{}`, and the comma-ok form of a lookup tells you whether a key is present.

For `Intersect`, start a new, empty `Set[T]{}`, range over `s`, and add each value that `other` has. Return the new set and leave `s` and `other` alone.
