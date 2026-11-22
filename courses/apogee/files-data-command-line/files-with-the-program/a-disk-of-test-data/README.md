---
slug: a-disk-of-test-data
title: A Disk of Test Data
type: code_tests
position: 5
difficulty: 7
topics: files, maps, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Mirela reviews the notice gctl prints when Cape Wren parks its dish safely in high wind, which the team calls a stow. Her rule: a test brings its own small disk and never reads the laptop's.

`fstest.MapFS` is a filesystem kept in a map. Each key is a slash-separated path, and each value an `*fstest.MapFile` whose `Data` holds the file's bytes. Folders named in a path exist by themselves, and a name the map lacks reports `fs.ErrNotExist`:

```go
logs := fstest.MapFS{
	"logs/night.txt": {Data: []byte("quiet\n")},
}
```

## Your task

Make `fixtureFS` in `fixtures.go` return an `fstest.MapFS` with two files, every line ending in `\n`. The first, `stations.csv`, holds the header and Cape Wren's line:

```
code,name,lat_deg,lon_deg,mask_deg,bands
CWR,Cape Wren,-53.9,-95.0,5,S X
```

The second, `templates/notices/stow.txt`, holds one line:

```
{station} stowed: {reason}
```
