---
slug: one-dish-down
title: One Dish Down
type: code_tests
position: 1
difficulty: 8
topics: errors, goroutines, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Tuesday 25 April: dishd's poll reported nothing at all while Dish 2's controller rebooted for a firmware check, not even Dish 1's good status.

Work that runs side by side should fail side by side: one dish's error must not erase the other's answer. A `Result` holds a dish's `Status` or its `Err`, so each poll can fail on its own, and the given `watchLine` shows a failed dish as `DOWN`, with the reason:

```go
for _, r := range results {
	fmt.Println(watchLine(r))
}
// YST-1 TRACKING az=183.25 el=41.70 wind=64
// YST-2 DOWN closed
```

## Your task

`pollAll` returns `nil` when any dish has failed. Return every dish's `Result` instead, each with its own `Status` or `Err`.

> [!TIP]
> `io.EOF` from a dish means its controller hung up between frames: `watchLine` calls that `closed`.
