Return at once for an empty `p`: a size of 0 would end the answer. Otherwise write the size line with `%x`, which prints lowercase hex, then the bytes, the CR LF after them, and `Flush`, whose error is the one to return: the writer keeps the first error of the writes before it.

A `for` loop with `range` over `lines` takes each line until the channel closes, which is how the store ends every feed. Each chunk is the line and a newline, as bytes.

If Go's client waits a second for the first line, a chunk is written but never flushed. If the feed ends with `unexpected EOF`, the zero chunk never comes, and if the collector that left gets nil, a `writeChunk` error is dropped.
