An `*os.File` is an `io.ReadSeeker`, so it can go to `http.ServeContent` as it is.

`file.Stat()` returns a `fs.FileInfo` whose `Name()` and `ModTime()` are the two values `ServeContent` wants beside the file; answer a `Stat` error through `writeError`.

If the resume that sends `If-Range` with the image's own date gets the whole file, the modification time you passed isn't the file's: a zero time or the clock's time never matches it.
