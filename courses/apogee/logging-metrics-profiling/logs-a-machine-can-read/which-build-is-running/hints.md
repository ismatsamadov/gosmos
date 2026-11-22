Call `debug.ReadBuildInfo()` at the top of `logStart`, after adding `"runtime/debug"` to the imports. It returns the build info and `ok`; read the info only when `ok` is true.

The three values hang off the build info: `GoVersion`, and `Main` with its `Path` and `Version`. Give each to `l.Info` after the message, a key followed by its value.

If the test says `"module" = "lesson.test"`, you logged `Path`, the path of the program's main package; the module's own is `Main.Path`.
