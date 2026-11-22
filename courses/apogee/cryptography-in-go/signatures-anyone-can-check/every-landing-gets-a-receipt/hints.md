Check before signing: a landing whose `BookingID` isn't the booking's `ID` is an error, and so is a digest `parseDigest` refuses (keep only its error, with `_, _, err :=`). Return nil and the error before any `Receipt` exists.

Fill the `Receipt` field by field: version 1, `n.KeyID`, the booking's `ID`, `SatID` and `Station`, the landing's `File`, `Bytes`, and `At` as `LandedAt`, and the digest as given.

`SignedAt` is `now()` without its fraction of a second: `.Truncate(time.Second)`. Then `wrapReceipt(n, r)` signs the receipt's JSON and returns the envelope. If `TestExactBytes` stops at `signed_at`, the milliseconds are still there.
