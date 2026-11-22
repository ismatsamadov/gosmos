---
slug: the-key-ring
title: The Key Ring
type: code_tests
position: 6
difficulty: 8
topics: crypto, errors, streams
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Thursday. At 00:00 UTC tomorrow `k-2027-03` retires, and from then on uplinkd derives every satellite's keys from the versions of the master key in its key file, `uplink-keys.txt`.

The file holds one version per line: its ID, the time it comes into use, the time it retires (`-` for not yet), its secret in hex, and the word `revoked` if it has been withdrawn. Lines starting with `#` and blank lines are notes for people. In August, below its notes, the file holds one version:

```text
k-2028-07 2028-08-11T00:00:00Z - 75706c69…6573
```

## Your task

1. Make `loadRing` pass every line that isn't blank or a note to `parseVersion`, and put the version in the ring with `ring.add`, which refuses an ID the ring already holds. Return any error with its line number first, as in `line 3: …`.
2. Make `SatKey` derive with `satKey` from the version named `keyID`, or wrap `ErrUnknownKey`: `unknown key: "k-2027-03"`.
