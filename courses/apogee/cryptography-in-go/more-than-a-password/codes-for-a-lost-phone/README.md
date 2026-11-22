---
slug: codes-for-a-lost-phone
title: Codes for a Lost Phone
type: code_tests
position: 7
difficulty: 9
topics: crypto, security, slices
created: 2026-11-13T06:18:24Z
modified: 2026-11-13T06:18:24Z
---
Thursday 2 November. Sunniva, from Ystmark: "Phone in a snowdrift. Snow 40 cm. Found in April, maybe." From Friday a password without the phone's code opens nothing, so every engineer also gets a sheet of eight recovery codes, each good for one sign-in in place of a code.

Each is `rand.Text()`, 130 bits, printed once and kept away from the phone. The console keeps only their digests, as it does for reset links, so a copy of its file signs nobody in. Random codes need no slow hashing: unlike passwords, they leave no short list of likely ones to try. A used code is spent, and `slices.Delete` removes an element:

```go
left := []string{"a", "b", "c"}
left = slices.Delete(left, 1, 2)
fmt.Println(left) // [a c]
```

## Your task

1. Make `Issue` keep only each code's `tokenDigest` in `r.unused`, and still return the codes.
2. Make `Use` accept `code` when its digest is among `user`'s unused ones, compared in constant time, and remove that digest.
