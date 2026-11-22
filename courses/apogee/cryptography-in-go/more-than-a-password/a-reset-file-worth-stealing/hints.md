`Start` already draws its token from `rand.Text`; only what it puts in `Hash` changes. The token itself goes back to the caller, for the link, and nowhere else.

`subtle.ConstantTimeCompare` takes two byte slices, so convert both hex strings with `[]byte(...)`. It returns 1 when they are equal, and a comparison with the number 1 is fine.

If the backup test still resets someone, `Redeem` compares the typed token with `Hash` itself: hash it first, so that a digest pasted from the file is hashed again and matches nothing.
