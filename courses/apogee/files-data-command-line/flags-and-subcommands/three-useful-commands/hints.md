Start as `runStations` does, with the stow synopsis, and define `-reason` with `fs.String` before parsing; its default is the reason in the second example line. After a successful parse, the flag set's `NArg` counts the arguments left after the flags, and `Arg` returns one by its position.

A wrong count is a usage error, so `usageError` writes the problem and the usage and gives you its status. An unknown code is a failure while working: build the error with `fmt.Errorf` and `%q`, as `runPasses` does, and hand it to `report` with no lines.

The record line is the code, then `stowed:` and the reason, one space apart. The flag's value is a pointer, so dereference it when you build the line, then pass the line to `report` in a slice of one.
