---
slug: three-parts-two-spaces
title: Three Parts, Two Spaces
type: code_tests
position: 1
difficulty: 8
topics: errors, http, strings
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
The recorder's request lines are tidy, but the station server will hear from anything on the hut's network. Nkem, from Port Ellery: "Accept one shape. Refuse the rest."

A request line is a method, a space, the target, a space and the version: `GET /recordings HTTP/1.1`. `strings.Split` keeps the empty part a doubled space leaves, which `strings.Fields` skips:

```go
fmt.Printf("%q\n", strings.Split("a  b", " "))
// ["a" "" "b"]
fmt.Printf("%q\n", strings.Fields("a  b"))
// ["a" "b"]
```

## Your task

`parseRequestLine` in `request.go` returns the whole line as the method. Require three non-empty parts split by single spaces, a target starting with `/` and a version `versionShaped` accepts, else return a `RequestError` 400 `bad request line`; any version but `HTTP/1.1` gets 505 `version "HTTP/1.0" not supported`.

> [!NOTE]
> Methods are case-sensitive: return `get` as sent, for a later check. Accepting only HTTP/1.1, and targets starting with `/`, is the station server's own strictness: HTTP also allows HTTP/1.0 and a full URL.
