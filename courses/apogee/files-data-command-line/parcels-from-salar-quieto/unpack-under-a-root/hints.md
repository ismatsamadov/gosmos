`filepath.Dir(f.Name)` is the folder a file sits in, `.` for a file at the top. `root.MkdirAll` makes it, and does nothing more when it already exists.

Create the file through `root`, never with `os.Create` and a joined path: only the root's methods refuse a link that leads out. Once it exists, close it on every way out, and after a good copy, return what `Close` returns.

Copy with `io.Copy` from `src` through `io.LimitReader` set to `limit+1`. `io.Copy` returns how many bytes it copied: more than `limit` means the file is too long, and exactly `limit` is fine.
