`res.Header.Get("Retry-After")` gives the header's text, or `""` when there is none, and `strconv.Atoi` turns whole-number text into an `int`, with an error for anything else, the empty string and dates included.

Decide on the number before making a duration: an error or a number below 1 means 1 second, a number above 60 means 60 seconds, and anything else is that many seconds.

If the case with 9999999999 seconds fails, you multiplied by `time.Second` before comparing with 60, and the product overflowed. Compare the plain `int` first, then multiply.
