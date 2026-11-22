---
slug: copy-the-recorder
title: Copy the Recorder
type: code_tests
position: 2
difficulty: 3
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-07T19:09:32Z
---
The van's next recording is larger than its quick preview. Paz wants to copy it without building one enormous string in memory.

`io.Copy(dst, src)` transfers bytes until the source ends or an error occurs. It returns the number written and any failure; normal end-of-file is success. The destination comes first, unlike many helper functions that take an input first:

```go
var saved bytes.Buffer
src := strings.NewReader("ready")
n, err := io.Copy(&saved, src)
fmt.Println(n, err) // 5 <nil>
```

## Your task

Implement `copyRecording` by returning `io.Copy`'s results. Do not discard the count when copying fails: the caller needs to know how much reached the destination before the failure.
