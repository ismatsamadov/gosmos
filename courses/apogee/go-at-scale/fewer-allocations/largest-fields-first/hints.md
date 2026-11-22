The dots in your picture were the gaps: six where 1-byte fields came before an 8-byte one, and two at the end, where the size was rounded up to a multiple of 8.

`unsafe.Alignof` is 8 for the two `int64` fields, 4 for the `float32` and 1 for the four 1-byte fields. Fields in order of their alignment, largest first, leave no gap between them.

If the test still says 32, look at where `MaxElevDeg` sits: the four 1-byte fields together fill exactly the 4 bytes beside it.
