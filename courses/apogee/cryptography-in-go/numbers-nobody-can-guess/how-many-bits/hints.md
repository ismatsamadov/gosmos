`math.Log2` takes a `float64`, so convert the alphabet's size with `float64(alphabet)` before passing it, and convert `length` the same way before you multiply.

Multiply only the length by the log: log2(10) is about 3.32, so six digits come to 6 × 3.32, about 19.9 bits.

If a line still prints 60.0 or 832.0, the product of `length` and `alphabet` is still in the `return`; the alphabet's size goes into `math.Log2` and nowhere else.
