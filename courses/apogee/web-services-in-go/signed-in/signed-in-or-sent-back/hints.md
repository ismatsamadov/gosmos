`r.Cookie(cookieName)` returns the cookie, or `http.ErrNoCookie` and a nil cookie. Check the error before you read the cookie's `Value`.

`a.Sessions.Lookup(cookie.Value)` gives the customer's ID, and `a.Customers[id]` the `Customer`. Each comes with its own `ok`, and both must be true.

If the case with another cookie before the session fails, read the cookie through `r.Cookie`, not the whole `Cookie` header: the browser puts every cookie it has for the service in that one header.
