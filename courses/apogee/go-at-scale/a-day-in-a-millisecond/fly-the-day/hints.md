`waitUntil(ctx, now, c.AOS)` returns `nil` at the contact's AOS, or at once when that AOS has passed, so the loop needs nothing else to keep time.

Call `fly(c)` only once the wait has returned `nil`. When it returns an error instead, return that error at once: the contacts after it stay unflown.

If the test says a contact flew at 00:00:00, `fly` is still called before any wait. If 993 contacts flew after the cancel at noon, the loop checks the context only between contacts: return the error `waitUntil` gives.
