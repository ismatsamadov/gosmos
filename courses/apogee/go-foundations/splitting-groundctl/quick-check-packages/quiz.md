## What is the import path of the code in `telemetry/`?

- [ ] `telemetry`
  > That's the name code uses after the import; the import path itself starts with the module path.
- [x] `apogee/groundctl/telemetry`
  > It's the module path from `go.mod`, then a slash and the folder.
- [ ] `./telemetry`
  > Go doesn't accept relative import paths in a module: you write the full path.

## Could `telemetry/limits.go` start with `package limits`?

- [ ] Yes: each file can be its own package
  > A package is a whole folder, not a file, and Go refuses a folder whose files name different packages.
- [ ] Yes, as long as the file is named after the package
  > File names don't decide anything here: every file in `telemetry/` says `package telemetry`.
- [x] No: a folder holds one package, so it says `package telemetry`
  > One folder is one package; code for a `limits` package would go in a folder of its own.

## Which function in `telemetry` can `main` call?

- [x] `Parse`
  > Its capital first letter exports it, so other packages can call it.
- [ ] `splitFields`
  > A lowercase name is unexported: only code inside package `telemetry` can call it.
- [ ] Both, once `main` imports `telemetry`
  > An import gives `main` only the exported names, the ones that start with a capital letter.
