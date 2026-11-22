---
slug: furrowlight-books-a-pass
title: Furrowlight Books a Pass
type: code_tests
position: 8
difficulty: 7
topics: errors, http, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Friday 3 December. Hollis's program is about to book Monday's 14:02 pass of Furrow-A over Salar Quieto, the first booking that will never touch Yusra's spreadsheet.

A create request runs its checks in order, and the first failure answers. When every check passes, the answer is `201 Created` with the new booking, and a `Location` header gives its address, so the caller never has to guess it:

```text
HTTP/1.1 201 Created
Content-Type: application/json
Location: /bookings/bk-4501
```

## Your task

`createBooking` decodes the body and keeps the booking. Finish it:

1. After `decodeJSON`, answer the first error of `checkShape`, then of `checkPass` with `s.Catalog` and `s.Now()`.
2. When `Create` returns an error wrapping `ErrConflict`, answer a 409 `conflict` `*APIError` with the error's own message.
3. Set `Location` to `/bookings/` and the new ID before the 201.

> [!TIP]
> Each check already returns the `*APIError` its failure deserves: hand it to `writeError` unchanged.
