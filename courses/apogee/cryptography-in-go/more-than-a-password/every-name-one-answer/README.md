---
slug: every-name-one-answer
title: Every Name, One Answer
type: code_tests
position: 2
difficulty: 8
topics: conditionals, security
created: 2026-11-13T06:18:24Z
modified: 2026-11-13T06:18:24Z
---
Thursday afternoon. Mirela types twenty names into the console's reset form. To some it answers `reset link sent`, to the rest `no account named` and the name. "It has read me our staff list."

The form must answer every name the same way, as the login has since August. `r.known` says whether an account exists, and `r.send` queues the email with the link; the answer is always `resetSent`. Mail goes only to a real account, whose owner then knows that someone asked:

```go
staff := map[string]bool{"mirela": true}
known := func(user string) bool {
	return staff[user]
}
fmt.Println(known("mirela"), known("root"))
// true false
```

## Your task

Make `Request` return `resetSent` for every name, and call `r.Start` and `r.send` only when `r.known(user)` says the account exists.

> [!NOTE]
> `send` only queues the email, so a real account's answer doesn't wait for the mail to go out.
