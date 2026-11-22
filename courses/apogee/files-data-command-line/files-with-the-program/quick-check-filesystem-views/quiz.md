## A function receives an fs.FS. Can it save a new report into it?

- [ ] Yes, with fs.WriteFile
  > The io/fs package has no WriteFile: its interfaces open and read, nothing else.
- [x] No, an fs.FS only opens and reads files
  > Saving a report takes the os package and a real path, such as a folder below t.TempDir in a test.
- [ ] Only when it is an embed.FS
  > Embedded files are part of the program itself, so they cannot be changed at all.

## After view, err := fs.Sub(resources, "templates"), which name opens templates/notices/stow.txt?

- [x] notices/stow.txt
  > The view's root is the templates folder, so names inside it leave that folder off.
- [ ] templates/notices/stow.txt
  > Through the view, that name looks for a templates folder inside templates, which does not exist.
- [ ] /notices/stow.txt
  > Names in a filesystem are relative: no fs.FS holds a name that starts with a slash.

## What is an fstest.MapFS for?

- [ ] Copying a folder from the laptop's disk into memory
  > A MapFS holds only what its map lists; nothing in it comes from the disk.
- [ ] Collecting the files that the code under test writes
  > Code can only read through an fs.FS; files a test checks after writing belong below t.TempDir.
- [x] Giving the code under test files the test defines in memory
  > Each key is a path and each value a file's bytes, so the test decides every file the code can see.
