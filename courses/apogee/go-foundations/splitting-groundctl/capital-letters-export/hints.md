The error is in `main.go`, but the cause is in `orbit/orbit.go`. What decides whether another package can use a name?

A name is exported when its first letter is a capital. Change the function's name where it's declared and where `main` calls it.

Go doc comments start with the name they describe, so the comment above the function changes with it. `cube` stays as it is, since only `orbit` calls it.
