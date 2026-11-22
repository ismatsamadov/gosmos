`errors.AsType[*APIError](err)` returns the `*APIError` it finds in `err` and whether it found one, looking through every error wrapped with `%w` along the way.

When it finds one, answer with `writeJSON`, the error's `Status`, and an `ErrorBody` holding its `Code` and `Message`, then return. The 500 that is there now stays, for every other error.

If the wrapped errors still get 500, you are checking the error's type directly instead of through `errors.AsType`. If the plain error's words reach the body, the 500 is built from the error instead of from the fixed words.
