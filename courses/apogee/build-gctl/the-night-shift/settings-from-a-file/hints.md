Open the file with `os.Open`; when `errors.Is(err, fs.ErrNotExist)`, return the defaults and `nil`, because a missing file is not a mistake. Return any other error. Add `"io/fs"` and `"os"` to the imports of `settings.go`.

Start from `defaultSettings()` in a variable and decode into its address: fields the file leaves out keep their defaults. Close the file when you are done, and wrap a decoding error with `fmt.Errorf`, the path and `%w`, so the night shift sees which file is broken.

In `apply`, convert the time to live with `time.Duration(st.CacheTTL)`. In `cmdPasses`, choose the code first, the typed argument or else `s.station`, then check and filter with that one code.
