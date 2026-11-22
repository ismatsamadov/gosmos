---
slug: open-where-you-left-off
title: Open Where You Left Off
type: code_tests
position: 5
difficulty: 8
topics: errors, files, json
created: 2026-11-07T06:18:32Z
modified: 2026-11-12T11:41:59Z
---
After Monday's deploy, the service counted from `bk-4501` again and handed out IDs that customers already had. After a restart it must know every booking, and never hand out `bk-4505` twice.

At start, a store reads back what it saved. No file means a first start, and an empty store. A file that doesn't decode means something went wrong: starting empty would let the next save write over every booking, so it must stop the start.

```go
data, err := os.ReadFile("gctl.json")
if errors.Is(err, fs.ErrNotExist) {
	return defaults, nil
}
if err != nil {
	return Settings{}, err
}
```

## Your task

`OpenFileStore` always starts empty. Make it open the file at the path it's given, a missing one meaning an empty store; decode it strictly into a `[]Booking`, refusing unknown fields and anything after the array; and restore the bookings with `f.mem.load`, which keeps their IDs and continues after the highest. A broken file returns an error and no store.
