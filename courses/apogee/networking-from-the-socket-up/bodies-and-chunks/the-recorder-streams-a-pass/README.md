---
slug: the-recorder-streams-a-pass
title: The Recorder Streams a Pass
type: code_tests
position: 8
difficulty: 9
topics: errors, http, networking
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Friday 9 June: ten booked Furrow passes stream into the station server today, and Hollis, at Furrowlight, has written in capitals what he expects: "EVERY JUNE PASS. EVERY FIELD."

An upload is kept only once everything about it holds: its path names a pass, its CRC matches, and the store doesn't have that pass yet. Then the answer says where the recording now is and what it holds:

```text
HTTP/1.1 201 Created
Location: /recordings/cirrine-1/YST/14140
Apogee-Crc32: 9f5edd58
```

## Your task

`put` in `routes.go` answers 201 to everything and keeps nothing:

1. Unless the target is `/recordings/` and a pass ID that `validPassID` accepts, answer 404 `no such path`.
2. Check the CRC with `checkCRC`; answer a refusal with `failure(err)`.
3. Keep the body with `s.Store.Put`, answering 409 with its error's text when the pass is already there. Then answer 201 `stored <pass ID> <n> bytes`, adding `Location` and `Apogee-Crc32`, in 8 hex digits.
