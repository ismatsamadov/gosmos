Read the header with `r.Header.Get` and give its value to `mime.ParseMediaType`, which returns the media type in lowercase without its parameters, and an error for a value it can't parse, an empty one included.

Do it first, before the decoder is made, and return the 415 when there is an error or when the media type is anything but `application/json`.

If `Application/JSON` is refused, you are comparing the header's whole value instead of the parsed media type. If a refused body was still read, the check comes after a `Decode`, or after the decoder has already started reading.
