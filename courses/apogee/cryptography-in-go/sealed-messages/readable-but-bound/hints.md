The fourth argument of `Seal`, and of `Open`, is the associated data. Both functions pass `nil` there now.

`s.AD()` returns the header as bytes, so it takes the place of the last `nil` in each call; the first `nil`, where the result goes, stays. In `sealWithHeader`, setting `s.Nonce` first changes nothing in `AD`, which reads only the four header fields.

If `TestHeaderBound` fails while the round trip passes, the two functions pass something other than the satellite's associated data: both must pass exactly `s.AD()`, not a header written some other way.
