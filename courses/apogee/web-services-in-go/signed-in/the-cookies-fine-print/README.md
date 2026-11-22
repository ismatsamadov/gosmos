---
slug: the-cookies-fine-print
title: The Cookie's Fine Print
type: code_tests
position: 0
difficulty: 8
topics: http, security
created: 2026-11-13T03:49:38Z
modified: 2026-11-13T03:49:38Z
---
Before a single cookie leaves the service, Mirela wants three answers: who can read it, where the browser sends it, and for how long.

A cookie is a name and a value that a server hands the browser in a `Set-Cookie` header. The browser keeps it and sends it back in a `Cookie` header with every request its attributes allow. `http.SetCookie` writes the header from an `http.Cookie`:

```go
http.SetCookie(w, &http.Cookie{
	Name:   "units",
	Value:  "km",
	MaxAge: 7 * 24 * 60 * 60, // a week
})
// Set-Cookie: units=km; Max-Age=604800
```

## Your task

`sessionCookie` returns a bare cookie. Give it these attributes:

1. `Path: "/account"`, so it goes back only to the account page, never to the API, and a `MaxAge` of `sessionLife`, in seconds.
2. `HttpOnly`, so no script on any page can read it, and `Secure`, so it travels only over HTTPS.
3. `SameSite: http.SameSiteLaxMode`, so the browser leaves it off forms that other sites post here.

> [!WARNING]
> `MaxAge` counts seconds: `30` would end the cookie after half a minute.
