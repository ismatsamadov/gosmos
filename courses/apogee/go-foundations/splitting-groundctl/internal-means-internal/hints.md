The error names the import that `main.go` may not use. Which package is allowed to import `keys`, and what does it offer instead?

Only code inside `uplink/` can import `uplink/internal/keys`. Import `uplink` itself, and call its exported function that returns the ID.

Change the import path and the value printed after `"signing with key"`; the `fmt` import stays.
