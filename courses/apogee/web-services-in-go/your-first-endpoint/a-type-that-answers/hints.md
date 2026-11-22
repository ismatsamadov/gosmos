The method's job is to put the list into `w`. `w` is an `io.Writer`, so anything that writes to a writer can write the answer.

`json.NewEncoder(w)` makes an encoder that writes to `w`, and its `Encode` method writes one value as JSON, followed by a newline. The value to encode is the receiver, `s`.

If the test says the body isn't JSON, something wrote Go's own printing of the list: encode with the JSON encoder rather than a `fmt` function.
