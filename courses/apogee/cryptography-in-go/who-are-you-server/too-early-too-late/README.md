---
slug: too-early-too-late
title: Too Early, Too Late
type: code_tests
position: 6
difficulty: 9
topics: crypto, time
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T07:46:48Z
---
Thursday 21 September. In Ingrid's rehearsal, the uplinkd at Lantern Atoll, the station near the equator, ran 40 minutes slow after a diesel restart and called a certificate 20 minutes old "not yet valid". Ingrid: "Exact times, please. When does one start counting, to the second?"

A certificate counts from `NotBefore` to `NotAfter`, both instants included: `Verify` refuses only a time before the first or after the last. It checks at `CurrentTime`, or at the Mac's clock when that is zero, and no test can move the Mac's clock. A clock passed in as a function can say any time:

```go
// a clock that always says 07:30 UTC:
now := func() time.Time {
	return time.Date(2028, 9, 21,
		7, 30, 0, 0, time.UTC)
}
fmt.Println(now().Format("15:04")) // 07:30
```

## Your task

`checkAt` checks at `time.Now()`. Make it check at the time `now()` gives, through `CurrentTime`. The test asks at `NotBefore` and `NotAfter` exactly, one second either side, and at Lantern Atoll's slow clock.
