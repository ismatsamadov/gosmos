---
slug: how-long-is-the-body
title: How Long Is the Body?
type: code_tests
position: 7
difficulty: 9
topics: conversions, errors, http
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
The recorder's test uploads say how long their bodies are, and on Tuesday Mirela's crooked client sent `Content-Length: 4096, 4096`: two programs could read that two ways.

Content-Length is one field holding one value: ASCII digits, no sign, no spaces. HTTP lets a server read a list of equal values as one; the station server refuses it. `strconv.Atoi` takes signs, so check the digits first, then parse them with `strconv.ParseInt`, which fails for a number too big:

```go
n, err := strconv.Atoi("+5")
fmt.Println(n, err) // 5 <nil>
```

## Your task

`contentLength` in `request.go` returns 0. Refuse, in order, two fields, anything but digits, more than `max` or too big to read (0 means no limit), and a body on GET or HEAD; a PUT with no field gets 411:

```text
400 more than one content-length
400 bad content-length "12a"
413 body over 65536 bytes
400 GET has no body
400 HEAD has no body
411 length required
```
