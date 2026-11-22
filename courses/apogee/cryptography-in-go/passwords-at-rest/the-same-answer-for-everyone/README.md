---
slug: the-same-answer-for-everyone
title: The Same Answer for Everyone
type: code_tests
position: 6
difficulty: 8
topics: errors, methods, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Thursday. Mirela tries the new login with a name nobody has, then with a real one: "no such user" comes back at once, "wrong password" 60 ms later. "Two answers, two timings. Anyone can list our 26 accounts."

A login must give an unknown user and a wrong password the same answer after the same work. When the name is unknown, the console checks the password anyway, against a dummy record: a valid 600,000-iteration record whose password nobody knows, so the check fails after one derivation, as a wrong password does. Today the console gives itself away:

```text
Check("intruder", …)  no such user     0 ms
Check("mirela", …)    wrong password  60 ms
```

## Your task

Make `Check` return `errInvalidLogin` for an unknown user and a wrong password alike. For a user who isn't in `c.records`, check the password against `dummyRecord` with `c.matches`, as if it were theirs, and refuse whatever it says.

> [!NOTE]
> `c.matches` derives through `c.kdf`, and the test's `kdf` counts its calls: every failed login must cost exactly one, at 600,000 iterations.
