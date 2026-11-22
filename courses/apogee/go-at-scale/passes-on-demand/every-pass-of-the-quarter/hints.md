A satellite's operator is `n.Satellites[i].Operator`, and its elements are `n.Elements[i]`, at the same index.

Build a `Network` that has all of `n.Sites` but only the operator's satellites and their elements, and hand it to `Catalog`: it will predict nothing else, still a day at a time.

The test fails if `predict` is asked about another operator's satellite, if a request with no operator calls it at all, or if a loop that stops at its twentieth pass has made it predict any day but 1 July.
