A `json.Decoder` with `DisallowUnknownFields` refuses a field your struct doesn't have, and after one `Decode`, the decoder's `Token` method returns `io.EOF` only when nothing follows the object. gctl read its settings file this way.

Put the default in your struct before you decode: a file without `days` keeps your 7, and one that says `"days":0` sets 0, which you refuse. `time.Parse` with the layout `time.DateOnly` accepts `2029-07-23` and refuses anything longer, a time of day included. `newNetwork().Sites` lists the twenty stations, in table order.

Wrap each problem with `fmt.Errorf` and `%w` after `rehearsal.json: `, and return a failed open's error as it is, so `errors.Is` still finds `fs.ErrNotExist` in it. Print the line only once every check has passed, with one `fmt.Fprintf` whose format ends in `\n`; for one, it's `day` and `station`, without the s.
