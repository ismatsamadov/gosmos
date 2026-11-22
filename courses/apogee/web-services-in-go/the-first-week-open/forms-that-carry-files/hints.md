Replace `r.Body` before anything reads it: `MultipartReader` reads from the body it finds when you call it, so a limit set afterwards limits nothing.

Call `NextPart` in a loop. Skip parts whose `FormName()` isn't `"file"`, read the one that is with `io.ReadAll`, and when `NextPart` says `io.EOF`, return `http.ErrMissingFile`.

If the 1.9 GB form ends in anything but a `*http.MaxBytesError`, something read the body without the limit: an `io.LimitReader` on the part stops quietly instead of saying the body was too large.
