`customerFrom(r.Context())` tells you whether `authenticate` left a customer. Without one, any error you build with `errors.New` makes `writeError` answer 500; return right after it.

`strconv.Itoa` turns `bucketSize` and the remaining count into header values. Set both headers before you decide, so a refused answer carries them too.

`Allow` has already rounded `retryAfter` up to whole seconds: divide it by `time.Second` and write the `int`, so the header reads `1`, never `1s`. On a refusal, set `Retry-After`, answer `errRateLimited` through `writeError` and return without calling `next`.
