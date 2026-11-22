Read the header's value with the header's `Get` method, which finds `Content-Type` whatever the letter case it was set in, and returns an empty string when it is missing.

`mime.ParseMediaType` returns three things: the media type in lowercase without its parameters, the parameters, and an error when the value can't be parsed, an empty value included. You need only the first and the last.

If `application/json; charset=utf-8` fails, you are comparing the whole value; if `application/jsonp` passes, you are comparing only its start. Compare the media type with `application/json` exactly, and only when the error is nil.
