The target is below `root`, not relative to the program's working directory. Build it with `filepath.Join`.

Use `MkdirAll` rather than `Mkdir`: the intermediate folders may not exist yet, and a second run must succeed.
