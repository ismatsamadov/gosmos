`http.NewCrossOriginProtection()` makes the check. Its `Handler` method wraps a handler, so that a refused request gets 403 and never reaches the handler inside.

Return the wrapped `h` from `noCrossOrigin`: every route of the account page is behind it then, the sign-in form's POST included.

GET and HEAD always pass, so a link from Yusra's email still opens the page. That is also why nothing on the page may change anything on a GET.
