---
slug: only-your-own
title: Only Your Own
type: code_tests
position: 8
difficulty: 9
topics: context, http, security, slices
created: 2026-11-07T06:18:33Z
modified: 2026-11-07T06:18:33Z
---
Mirela's headline for Friday: "Brinewatch must never see Cirrine's passes." Today any customer with a token still lists every satellite's passes and bookings.

Knowing who is asking is authentication; deciding what they may see and do is authorization. Another customer's pass or booking answers 404, exactly like a missing one, and booking another customer's satellite answers 403 before the catalog is asked, so nobody learns what others have:

```go
bands := []Band{"UHF", "S"}
fmt.Println(slices.Contains(bands, "X"))
// false
```

## Your task

1. In `auth.go`, write `owns`: true when `c` operates `satID`.
2. In `passes.go` and `bookings.go`, take the caller with `customerFrom`. Keep `listPasses` and `listBookings` to its satellites, and make `getPass`, `getBooking` and `cancelBooking` answer for another customer's exactly as for a missing one.
3. In `createBooking`, between the `checkShape` and `checkPass` checks, refuse another customer's satellite with 403 `forbidden` and a message such as `cirrine-1 is not one of your satellites`, naming the booking's `sat_id`.

> [!WARNING]
> After `checkPass`, a missing Cirrine-1 pass would get 422 and a real one 403, telling Brinewatch which passes Cirrine-1 has.
