---
slug: walk-the-templates
title: Walk the Templates
type: code_tests
position: 4
difficulty: 7
topics: files, functions, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Before the climb, Paz wants gctl to list every template it carries, the nested ones too, so the night shift at Salar Quieto knows which reports they can start from.

`fs.WalkDir(fsys, root, fn)` calls `fn` for `root` and for every file and folder below it, each folder before what it holds, in lexical order. `fn` receives a path, an `fs.DirEntry` and an error. When that error is not nil the entry may be nil, and returning the error stops the walk:

```go
err := fs.WalkDir(fsys, "logs",
	func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		fmt.Println(p, d.IsDir())
		return nil
	})
```

## Your task

Implement `listTemplates`: walk `templates` in `fsys` and return the path of every file in walk order, leaving folders out. Return the walk's error, such as the one for a missing `templates` folder.
