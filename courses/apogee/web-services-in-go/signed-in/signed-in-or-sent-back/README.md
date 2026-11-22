---
slug: signed-in-or-sent-back
title: Signed In or Sent Back
type: code_tests
position: 2
difficulty: 8
topics: errors, http, maps
created: 2026-11-13T03:49:38Z
modified: 2026-11-13T12:18:55Z
---
Yusra signs in with Apogee's own key to try the page on her test booking, Pipit-1 over Port Ellery on Monday, since Apogee tries everything on a Pipit first. The page sends her straight back to the sign-in form.

`page` asks `signedIn` who the browser is, and `signedIn` doesn't look. A browser sends all its cookies for the service in one `Cookie` header, and `r.Cookie` picks one out by name. When there is none of that name, it returns `http.ErrNoCookie`: nothing went wrong on the service; the browser isn't signed in.

```go
c, err := r.Cookie("units")
if errors.Is(err, http.ErrNoCookie) {
	return "km" // nothing chosen yet
}
return c.Value
```

## Your task

`signedIn` says nobody is signed in. Read the cookie named `cookieName`; without one, return false. Look its value up in `a.Sessions`, find that customer in `a.Customers`, and return the cookie's value, the customer and true. If the session or the customer isn't found, return false.
