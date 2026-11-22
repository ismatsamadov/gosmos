---
slug: headers-name-the-body
title: Headers Name the Body
type: code_tests
position: 6
difficulty: 5
topics: http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra's TL;DR from Cirrine: "Their test program threw away a good answer." Its `Content-Type` said `application/json; charset=utf-8` and their check compared the whole line, so Yusra wants the customer guide to show a check that works.

`Content-Type` names the body's format: a media type such as `application/json`, sometimes followed by parameters after a `;`. `Header.Get` finds a header whatever the letter case of its name, and `mime.ParseMediaType` splits a value, returning the media type in lowercase:

```go
mt, params, err := mime.ParseMediaType(
	"Text/HTML; charset=utf-8")
fmt.Println(mt, params["charset"], err)
// text/html utf-8 <nil>
```

## Your task

`isJSON` returns `false`. Make it return true when the media type in `h`'s `Content-Type` is `application/json`, with any parameters and in any letter case, and false when the header is missing or can't be parsed.

> [!TIP]
> Compare the whole media type: `application/jsonp` starts with `application/json`, but it isn't JSON.
