After the `\n` is gone, the line may still end in a `\r`, the one the controller sends and page one leaves out. `strings.TrimSuffix` removes a suffix only when it is there.

Remove one `\r` from the end and no other byte, so `YST-1 READY PROTO 1\r\n` and page one's `YST-1 READY PROTO 1\n` both leave `YST-1 READY PROTO 1` for the checks below.

If the line with two carriage returns or a trailing space is accepted, something removes more than one `\r`: `TrimSpace` and `TrimRight` remove every space or `\r` they find at the end.
