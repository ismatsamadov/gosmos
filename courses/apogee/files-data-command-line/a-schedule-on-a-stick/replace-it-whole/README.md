---
slug: replace-it-whole
title: Replace It Whole
type: code_tests
position: 7
difficulty: 3
topics: files, testing
created: 2026-11-04T17:28:41Z
modified: 2026-11-05T08:00:06Z
---
A partial save must not replace the van's last usable schedule. Paz wants readers to see the old file until the new one is complete.

Create a temporary file in the destination's directory, write it, close it and rename it over the destination. Keeping both files in one directory avoids a rename across filesystems. Remove the temporary name on every failure. This pattern controls what readers see; it does not promise survival after a power failure.

```go
file, err := os.CreateTemp("reports", ".draft-*")
if err != nil {
    return err
}
defer os.Remove(file.Name())
defer file.Close()
```

## Your task

Implement `replaceSchedule` with `CreateTemp`, `Write`, `Close` and `Rename`, checking each error. Use `filepath.Dir(path)` for the temporary file's directory. Close it before renaming, and arrange cleanup even when writing or renaming fails.
