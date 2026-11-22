---
slug: finding-aos
title: Finding AOS
type: code_tests
position: 4
difficulty: 4
topics: loops
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
AOS, the start of a pass, is the first minute Pipit-3 stands at or above the mask. Once you've found it, there's no need to read the rest of the prediction.

`break` leaves a loop at once; the program carries on after the loop's closing brace:

```go
for try := 1; try <= 5; try++ {
	if echoHeard(try) {
		fmt.Println("echo on try", try)
		break
	}
}
```

## Your task

Make `aosMinute` look through minutes 0 to 14 and return the first whose `elevationAt` is at or above `maskDeg`, with `break` stopping the search there. If no minute qualifies, return -1.
