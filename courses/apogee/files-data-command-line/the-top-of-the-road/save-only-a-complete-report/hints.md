Make the temporary file with `os.CreateTemp` in the destination's folder, which `filepath.Dir` gives, with a pattern such as `.day-*`. Right after it succeeds, defer removing its name and closing it, so every early return cleans up.

Hand the temporary file to `report` and return its error before anything else happens: the destination has not been touched yet. Then check the error from `Close`, and only after that rename the temporary file's name to `path`.

After a successful rename, the deferred removal finds nothing under the temporary name, and the deferred close of an already closed file does no harm.
