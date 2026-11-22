---
slug: signing-out-for-real
title: Signing Out for Real
type: code_tests
position: 3
difficulty: 8
topics: http, security
created: 2026-11-13T03:49:38Z
modified: 2026-11-13T12:18:56Z
---
Mirela clicks Sign out, puts her old cookie back into the browser by hand, and reloads: Apogee's bookings are back on the screen. Signing out only showed her the sign-in form.

Signing out has two halves. The service forgets the session, so its ID signs nobody in again, even from a copy. And the browser is told to drop its cookie, with a `Set-Cookie` of the same name and path and `Max-Age=0`, which `MaxAge: -1` writes. A cookie with another path is another cookie, and the browser would keep the old one.

```go
c := &http.Cookie{Name: "units", Path: "/guide"}
c.MaxAge = -1
http.SetCookie(w, c)
// Set-Cookie: units=; Path=/guide; Max-Age=0
```

## Your task

`signOut` only redirects. Before it does, read the session cookie and, if there is one, end its session with `a.Sessions.End`. Then send the browser `sessionCookie("")` with its `MaxAge` set to -1.
