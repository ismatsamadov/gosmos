Start with `c.Check(user, pw)`: when it returns an error, return an empty token and that error at once, before anything is saved or added.

For the upgrade, pass `c.records[user]` to `needsUpgrade`. When it says yes, `c.newRecord` makes the new record, `c.save` stores it, and an assignment to `c.records[user]` lets the next login find it. Return any error from either call.

For the token, `sha256.Sum256` of the token's bytes and `hex.EncodeToString` of the result give the key for `c.sessions`. If the kiosk is saved twice, `c.records` still holds the old record when the second login checks it.
