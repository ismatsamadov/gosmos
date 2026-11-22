Follow a request through the list: the first middleware listed sees it first and sees the answer last. For each one, ask what must already have happened when it runs, and what it needs to see afterward.

`logRequests` reads the ID from the request it is given, so `requestID` has to come before it. It writes its line after the handlers inside it return, so the 500 that `recoverPanics` writes must be written inside it.

If the panicking request's access line is still missing, `recoverPanics` sits outside `logRequests`, and the panic unwinds through `logRequests` before it can write. `noStore` can go last, right around the routes.
