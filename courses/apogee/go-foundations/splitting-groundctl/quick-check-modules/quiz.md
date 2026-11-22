## What does the `module` line set?

- [ ] The folder the module must be stored in
  > A module can sit in any folder on your disk; the module line names it, wherever it lives.
- [x] The path that every package's import path starts with
  > `apogee/groundctl/orbit` is that path plus a folder, and so is every other package in the module.
- [ ] The Go version the module needs
  > That's the `go` line below it.

## What does `go.sum` hold?

- [ ] The version of each dependency to build with
  > The `require` lines in `go.mod` choose the versions; `go.sum` checks what was downloaded for them.
- [ ] A list of the module's own packages
  > Go finds your packages from their folders; no file lists them.
- [x] A checksum of each dependency's contents
  > The go command compares every download with it, so a dependency that changed under the same version is refused.

## groundctl v2.0.0 breaks its API. What is the import path of its `orbit` package?

- [x] `apogee/groundctl/v2/orbit`
  > From version 2 on, the module path ends in the major version, here `/v2`, and every package's import path starts with the module path, so a program can even use v1 and v2 side by side.
- [ ] `apogee/groundctl/orbit`, with v2.0.0 in the `require` line
  > That path stays with v1: a new major version changes the module path itself, not only the version in `go.mod`.
- [ ] `apogee/groundctl/orbit@v2`
  > `@` and a version belong on the command line, as in `go get`, never in an import path.
