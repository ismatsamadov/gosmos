---
slug: fixtures-in-testdata
title: Fixtures in testdata
type: code_tests
position: 6
difficulty: 6
topics: testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Thistle's recent frames become fixtures: twelve fixed input lines that the parser's tests read from a file.

`go test` runs in the package's folder, and the `go` tool ignores folders named `testdata`, so files for tests go there. A test that can't read its input has nothing left to check, so it stops with `t.Fatalf`:

```go
id, err := parsePassID("pipit-1/PEL/04182")
if err != nil {
	// nothing below makes sense without id
	t.Fatalf("parsePassID: %v", err)
}
```

## Your task

Write `loadFixture` in `fixture.go`: read `testdata/<name>` with `os.ReadFile`, which returns the file's bytes and an error. If that fails, stop the test with `t.Fatalf` and a message naming the file. Otherwise return the file's non-empty lines; `string(b)` turns a byte slice `b` into text.
