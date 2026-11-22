## What does filepath.Join do when the resulting path does not exist?

- [ ] It creates the missing directories
  > Directory creation is a separate operation such as MkdirAll.
- [ ] It returns a not-exist error
  > Join operates on strings and does not inspect the disk.
- [x] It returns the cleaned path string
  > A path can be constructed before any of its components exist.

## What happens when MkdirAll reaches an existing directory?

- [x] It accepts it and continues
  > This makes repeated preparation of a folder safe.
- [ ] It deletes the directory first
  > MkdirAll does not remove an existing directory or its contents.
- [ ] It always returns an error
  > An existing directory is allowed; a blocking regular file is not.

## What does FileMode.Perm select?

- [ ] The filename's extension
  > A file's name is separate from its mode.
- [x] The nine owner, group and other permission bits
  > Perm excludes the additional bits used for file type and other mode information.
- [ ] The file's contents
  > Reading contents requires a read operation, not a mode inspection.
