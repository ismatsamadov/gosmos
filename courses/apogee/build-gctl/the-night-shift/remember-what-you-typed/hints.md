In `Run`, call `s.history.add` with the scanned line inside the loop, before `dispatch`, so a command that fails, or `history` itself, is recorded too.

`add` checks the line with `strings.TrimSpace`: only a line with something left is appended to `h.lines`, unchanged, spaces and all. Add `"strings"` to the imports of `history.go`.

`cmdHistory` ranges over `s.history.lines` with the index, and prints the index plus 1 through `%3d`, two spaces, then the line. An argument returns the error `history takes no arguments`; add `"errors"` and `"fmt"` to the imports.
