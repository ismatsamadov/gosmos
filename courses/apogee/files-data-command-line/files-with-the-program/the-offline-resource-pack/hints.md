Read `calibration.json` from `override` first, with `fs.ReadFile`. Only when that error is `fs.ErrNotExist`, as `errors.Is` tells, read the same name from `fsys` instead.

After that single fallback, any error left, from either read, is returned. Then parse the bytes with `parseOffsets` and return its error too: a broken file on the stick is not a missing one.

If the test says the cut-off file or the unreadable stick fell back to zero, the code falls back on more than a missing file: only `fs.ErrNotExist` from `override` may lead to the default, and a file that does not parse never does.
