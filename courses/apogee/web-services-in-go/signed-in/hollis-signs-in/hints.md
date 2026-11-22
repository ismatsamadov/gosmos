`chain(h, a, b, c)` runs `a` first. A refusal carries only the headers that the middleware outside the refusing one has already set, so `securityHeaders` and `noStore` come before `noCrossOrigin` in the list.

`cors(origin)` returns a middleware: call that on `api`, and register the result at `/` in place of `api`.

If the planning page's `POST /bookings` gets 403, `noCrossOrigin` is around the API too; if a page under `/account/` answers with `Access-Control-Allow-Origin`, `cors` is around the account page.
