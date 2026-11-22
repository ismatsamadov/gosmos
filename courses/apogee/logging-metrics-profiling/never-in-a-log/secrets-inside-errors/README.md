---
slug: secrets-inside-errors
title: Secrets Inside Errors
type: code_tests
position: 6
difficulty: 8
topics: errors, logging, security, strings
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
The 14th line Mirela found was a `request failed` record whose error quoted what Furrowlight's program had sent, the API key itself, inside the text under `err`, where a rule about attribute names never looks.

A secret inside free text needs a pattern. The `regexp` package finds text by a regular expression, such as `[0-9a-f]{32}` for 32 hex digits; a part in parentheses is a group, which a replacement names as `$1`. `secretPattern` matches Apogee's API keys and tokens, with the customer ID as its first group. Compile a pattern once, outside the function that uses it: compiling it inside redoes the same work for every string of every record.

```go
var orb = regexp.MustCompile(`/([A-Z]{3})/\d{5}`)

func noOrbit(id string) string {
	return orb.ReplaceAllString(id, "/$1")
}
// noOrbit("brine-3/LAN/14384") is "brine-3/LAN"
```

## Your task

Compile `secretPattern` once, into a variable outside `scrub`, and make `scrub` replace every match with `$1.[redacted]`, which keeps the customer ID and drops the secret.

> [!NOTE]
> The pattern finds keys and tokens in the shapes Apogee issues; a secret of any other shape, such as a password typed into a message, needs a rule of its own.
