The package-level `flag.String` and `flag.Parse` read the whole process's arguments. Make a flag set of the command's own with `flag.NewFlagSet`, and choose where it writes with its `SetOutput` method.

Define `-station` on that flag set with `StringVar`, pointing at the options' `Station` field, then parse `args` with the flag set's `Parse`. When it fails, return empty options and its error.

When a mistake's message is missing from `output`, the flag set still writes to the process's standard error: call `SetOutput` before `Parse`.
