Build the `Settings` value with its three defaults before decoding anything, and pass its address to `Decode`: the file then overwrites only the fields it names.

Call `DisallowUnknownFields` on the decoder before the first `Decode`. The second `Decode` can fill any target, such as a `json.RawMessage`; any result other than `io.EOF` means something follows the object.

Check the day last, by parsing it with `time.Parse` and the layout `time.DateOnly`. An empty day fails to parse too, so that one check covers a missing day and a wrong one.
