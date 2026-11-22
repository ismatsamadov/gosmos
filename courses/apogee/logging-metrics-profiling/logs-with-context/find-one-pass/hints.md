Make `bufio.NewScanner(r)` and loop with `for sc.Scan()`; inside, decode `sc.Bytes()` into a fresh `Line` with `json.Unmarshal`. Remove the TODO line once the imports are used for real.

When the decode returns an error, add one to `broken` and `continue`; otherwise keep the line only if its `PassID` equals `passID`, after setting its `Raw` to `sc.Text()`, by appending it to `lines`. After the loop, return `lines`, `broken` and `sc.Err()`.

If the test finds every line's `Raw` empty, it is set after the line was appended, on a copy; if the reader's error is missing, `sc.Err()` isn't what `follow` returns.
