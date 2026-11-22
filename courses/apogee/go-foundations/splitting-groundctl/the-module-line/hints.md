The compiler doesn't know what `orbit` means in `main.go`. Which part of a Go file tells it where a package comes from?

An import path is the module path from the `module` line in `go.mod`, a slash, then the package's folder.

With two imports, group them in one block: `import (`, then `"fmt"` and `"<module path>/orbit"` on lines of their own, then `)`.
