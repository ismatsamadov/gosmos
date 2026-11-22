---
slug: the-recorder-pipeline
title: The Recorder Pipeline
type: code_tests
position: 8
difficulty: 4
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-04T17:36:06Z
---
Paz wants one bounded diagnostic recording on the screen and in the archive. The van must consume it once, keep both destinations in step and report a failure instead of quietly throwing the recording away.

The stream helpers compose: put the limit around the source, combine the destinations, then copy. Validate the settings before consuming any bytes. A failed destination can leave earlier destinations with a partial recording, so keep the returned error.

```go
limited := io.LimitReader(source, 128)
dest := io.MultiWriter(&screen, &archive)
n, err := io.Copy(dest, limited)
```

## Your task

Implement `recordPreview` for a nonnegative `limit` and at least one target. Return an error for invalid settings without reading `src`. Otherwise copy at most `limit` bytes to every target and return the copy count and error. Zero is a valid empty preview.
