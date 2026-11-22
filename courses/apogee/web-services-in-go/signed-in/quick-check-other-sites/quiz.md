## Why does the account page need a guard against other sites' forms when the API doesn't?

- [ ] The API checks a signature, and the account page checks nothing
  > The page checks its session as carefully as the API checks a token; the difference is in how each reaches the service.
- [x] A browser adds the cookie by itself, but a token goes only where a page's code puts it
  > Unless `SameSite` keeps it off, the browser adds Hollis's cookie to a form from any page; no other page can add his token, which it never had.
- [ ] A form can only post to addresses under `/account/`
  > A form can post to any address, `/bookings` included; without a token the API refuses it there.

## A page on another site posts a form to the account page. What does CORS do about it?

- [ ] The browser never sends the form without `Access-Control-Allow-Origin`
  > A browser sends a form wherever it points. CORS is read from the answer, after the request has arrived.
- [ ] The browser asks with a preflight first, and the service refuses it
  > A plain form posts with no preflight; preflights come before requests that scripts send with headers such as `Authorization`.
- [x] Nothing: CORS decides which pages' scripts may read answers, not which requests arrive
  > That is why the page needs `noCrossOrigin`: only the service's own check stops the request itself.

## Why does `cors` add `Vary: Origin` to every answer?

- [x] The answer depends on the request's `Origin`, so a cache must never hand one origin's copy to another
  > The planning page's answer carries `Access-Control-Allow-Origin`; a copy kept for it must not reach any other origin.
- [ ] Browsers refuse an answer to another origin without it
  > Browsers read `Access-Control-Allow-Origin` for that. `Vary` speaks to caches.
- [ ] It names the one origin the service trusts
  > That is `Access-Control-Allow-Origin`'s job; `Vary` lists the request headers an answer depends on.
