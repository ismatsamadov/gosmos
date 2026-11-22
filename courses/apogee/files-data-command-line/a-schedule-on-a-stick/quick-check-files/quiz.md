## What does os.ReadFile return for a path that does not exist?

- [ ] An empty schedule and no error
  > Empty contents can be valid, but a missing path produces an error.
- [x] A byte slice and a non-nil error
  > The error explains why no file contents could be read.
- [ ] It creates the file first
  > ReadFile reads an existing path; it does not create files.

## How do you turn the successful bytes into text without trimming them?

- [x] string(data)
  > Conversion preserves all the bytes, including spaces and line endings.
- [ ] strings.TrimSpace(string(data))
  > Trimming removes bytes from the edges, so the text is no longer exact.
- [ ] fmt.Sprintln(data)
  > That formats the slice rather than converting its contents to text.

## Does WriteFile(path, data, 0600) change the permissions of an existing file?

- [ ] Yes, every write resets them
  > The creation mode is used when a new file is created.
- [ ] Only when the new contents are shorter
  > Contents and permission bits are separate properties.
- [x] No, its existing permissions remain
  > WriteFile replaces contents but does not change an existing file's permissions.
