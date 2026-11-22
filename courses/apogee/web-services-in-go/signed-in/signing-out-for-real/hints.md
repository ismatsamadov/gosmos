Read the cookie as `signedIn` does, with `r.Cookie(cookieName)`, and call `a.Sessions.End(cookie.Value)` only when that gave no error: a browser without a cookie has nothing to end.

`sessionCookie("")` already has the right name and path. Set its `MaxAge` to -1 and hand it to `http.SetCookie` before `http.Redirect`, which writes the answer's headers.

If the old cookie still opens the page, the service never forgot the session: clearing the browser's copy alone leaves the ID working wherever it was copied.
