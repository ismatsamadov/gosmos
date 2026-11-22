After the first `Decode`, ask the same decoder for one more token. Its error at the end of the input is the signal you want.

Compare that error with `io.EOF`. Anything else, a token with no error or a syntax error, means the file holds more than one object.

`dec.More()` looks tempting, but it reports false for a stray `}`, so the file would pass. The token check refuses it.
