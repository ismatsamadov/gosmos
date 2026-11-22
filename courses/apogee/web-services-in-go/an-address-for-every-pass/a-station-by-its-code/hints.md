The new pattern goes beside `GET /stations` in `routes`, with `s.getStation` as its handler. The word in the braces is the name you will ask the request for later.

The request's `PathValue` method takes a wildcard's name and returns the segment it matched, exactly as the caller wrote it. Use what it returns in place of the empty string.

If `/stations/SQT` gets `404 page not found`, the route isn't registered, or its wildcard has another name. If it gets `no station ""`, `getStation` still uses the empty string.
