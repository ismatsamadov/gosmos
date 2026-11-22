---
slug: tokens-people-can-type
title: Tokens People Can Type
type: code_tests
position: 3
difficulty: 6
topics: crypto, strings
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
When engineers log in to the web tool that runs Apogee's stations, the ops console, their browser keeps a code that stands in for the password until they log out: a session token. Mirela finds that each token is the login time in nanoseconds: "Know when Ingrid logged in, and you are Ingrid."

`rand.Text` from `crypto/rand` returns 26 characters from an alphabet of 32, `A` to `Z` and `2` to `7`, with no 0 or 1 to mistake for O or I. That is 130 bits, short enough to paste, with nothing in it about you or the time. A token made from the clock gives both away:

```go
tok := "1846972800123456789"
ns, _ := strconv.ParseInt(tok, 10, 64)
fmt.Println(time.Unix(0, ns).UTC())
// 2028-07-12 00:00:00.123456789 +0000 UTC
```

## Your task

`newSessionToken` returns the time in nanoseconds as text. Make it return `rand.Text()` instead. The test plays someone who knows, to the second, when you logged in.
