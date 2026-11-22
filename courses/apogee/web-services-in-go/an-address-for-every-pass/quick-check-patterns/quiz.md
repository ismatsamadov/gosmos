## What does `HEAD /stations` get from these two routes?

- [ ] `405 Method Not Allowed`, since only GET is registered
  > A `GET` pattern also answers HEAD, which asks for the same answer without its body.
- [x] `200`, with the list's headers and no body
  > The mux serves HEAD with the GET route's handler, and the server leaves the body out.
- [ ] `404 page not found`, since no pattern names HEAD
  > The path matches a route, and GET's pattern covers HEAD; 404 is for a path no pattern matches.

## What do `POST /stations/SQT` and `GET /stations/SQT/x` get?

- [ ] `404` for both
  > `/stations/SQT` matches a route under GET, so a POST there gets 405, not 404.
- [ ] `405` for both
  > `/stations/SQT/x` has one segment too many for `{code}`, so no route matches it at all.
- [x] `405` for the POST and `404` for the GET
  > The POST's path matches a route under another method; the GET's path matches no route.

## What does `{code}` match in the pattern `GET /stations/{code}`?

- [x] Exactly one segment, which can't be empty
  > `/stations/SQT` matches, while `/stations/` and `/stations/SQT/x` don't.
- [ ] One segment, which may be empty
  > A wildcard needs at least one character, so `/stations/` gets 404.
- [ ] Everything after `/stations/`, slashes included
  > A wildcard stops at the next slash; only a wildcard ending in `...` takes the rest of the path.
