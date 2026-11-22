---
slug: read-from-the-pipe
title: Read from the Pipe
type: code_tests
position: 5
difficulty: 6
topics: files, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
At this morning's stop, Paz connected the recorder to gctl through a pipe, which feeds one program's output into another program's input. No file is involved: the bookings arrive on standard input.

Command-line tools read standard input when no file is named, so one command works on a file from the stick and on the recorder's pipe. Open the file only when there is a name, and close it when you are done. Standard input belongs to whoever started gctl, so leave it open.

```go
table := defaultTable
if path != "" {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	table = string(data)
}
```

## Your task

Implement `loadBookings`: with an empty `name`, read the bookings from `stdin` with `readBookings`; otherwise open the named file, read its bookings and close it. Return an open error as it is.
