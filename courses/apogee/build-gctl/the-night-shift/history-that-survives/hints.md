`load` returns `nil` at once for an empty `h.path`, and when `os.ReadFile` fails with `fs.ErrNotExist`. Otherwise it adds each line of the file to `h.lines` without its newline; `strings.Lines` gives the lines one at a time, each still ending in its newline.

A small helper that appends one line keeps `add` short: open with the three flags, write the line with `fmt.Fprintln`, close the file, and return what went wrong in either step: `errors.Join` of the write's error and `Close`'s is `nil` only when both are. `add` calls it only when `h.path` isn't empty.

Set `h.failed` the first time the file fails, in `load` or in `add`, and skip the file while it is set, so the error is printed once. `add` now returns an `error`: in `Run`, print it with `error:` in front, then dispatch the line as usual. Add `"io/fs"` and `"os"` to the imports of `history.go`.
