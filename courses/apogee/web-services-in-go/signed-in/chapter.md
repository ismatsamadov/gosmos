---
slug: signed-in
title: Signed In
position: 10
---
Browsers: a session cookie with `HttpOnly`, `Secure`, `SameSite`, `Path` and `MaxAge`, a server-side session behind it with an ID from `crypto/rand.Text` that ends 30 minutes after its last use, reading it back with `r.Cookie`, signing out on both sides, `http.CrossOriginProtection` against other sites' forms, CORS for one trusted origin, and the headers that tell a browser how to treat a page.
