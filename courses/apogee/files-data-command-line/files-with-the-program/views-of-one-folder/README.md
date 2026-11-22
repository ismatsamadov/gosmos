---
slug: views-of-one-folder
title: Views of One Folder
type: code_tests
position: 6
difficulty: 7
topics: files, interfaces, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:21:59Z
---
The report renderer asks for `day.txt`, but among gctl's built-in files it is `templates/day.txt`. Nkem would rather hand the renderer a view of that one folder than teach it the prefix.

`fs.Sub(fsys, dir)` returns a filesystem whose root is `dir`: names in the view leave the folder off and cannot reach outside it. `Sub` does not check that the folder exists, so ask `fs.Stat` first when a missing folder is a mistake:

```go
manuals, err := fs.Sub(fsys, "docs/manuals")
if err != nil {
	return err
}
data, err := fs.ReadFile(manuals, "stow.txt")
```

## Your task

Implement `templateFS`: return the error from `fs.Stat` if `fsys` has no `templates` folder, and otherwise both results of `fs.Sub` for that folder.
