Start with an empty slice of `T`, look at every item, and keep the ones `keep` approves. Which loop visits each item of a slice in order?

A `for` loop with `range` over `items` hands you each item in turn. Call `keep` on it, and `append` the item to your result when `keep` says yes.

The shape: declare an empty `[]T`, loop with `range`, append when `keep(____)` is true, and return the slice you built.
