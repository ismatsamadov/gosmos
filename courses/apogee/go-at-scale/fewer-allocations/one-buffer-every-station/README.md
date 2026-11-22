---
slug: one-buffer-every-station
title: One Buffer for Every Station
type: code_tests
position: 4
difficulty: 9
topics: performance, slices, streams
created: 2026-11-12T03:10:33Z
modified: 2026-11-12T03:10:33Z
---
Wednesday 23 May. The export to the dish loaders makes a new 8 KB buffer for every station's day and throws it away after one write: twenty buffers for take-over day alone.

One buffer can serve every station in turn. `buf = buf[:0]` keeps the array and sets the length to 0, so the next station's lines overwrite the last one's. That's safe because `io.Writer` promises that `Write` doesn't keep the slice it's given: once `Write` returns, the bytes are yours again.

```go
line := make([]byte, 0, 64)
for _, f := range frames {
	line = line[:0]
	line = append(line, f.SatID...)
	line = append(line, '\n')
	os.Stdout.Write(line)
}
```

## Your task

Have `exportStations` make its buffer once, before the loop, empty it before each station's lines, and write it as before. Take-over day's export should then take at most 2 allocations instead of 20.

> [!NOTE]
> Ystmark's day is the longest, 7,811 bytes, so the 8 KB buffer never has to grow.
