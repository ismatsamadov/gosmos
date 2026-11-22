`DisallowUnknownFields` is a method of `json.Decoder`, not an option of `json.Unmarshal`: make a decoder for `r` first.

Call it before `Decode`. The decoder checks each field name while it decodes, so switching the check on afterwards changes nothing.

When `Decode` fails, return `Settings{}` rather than the value it partly filled: the fields that did match must not be used from a refused file.
