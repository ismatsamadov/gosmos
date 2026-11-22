Start `logResult` with an `if` on `err != nil`: that branch writes the failure and returns, and the `l.Info` line after it stays for a flown pass.

In the failure branch, write the record with `l.Error`: the message first, then `pass_id` and `err`, each key followed by its value. Give it the error itself, and slog writes its whole text.

If the flown pass fails with `"err":null`, the `err` pair is written on every call; it belongs only in the branch where `err` isn't nil. If `"err"` holds only `context deadline exceeded`, the error was unwrapped before it was logged.
