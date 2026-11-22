`p.Expires` is the first instant the link is refused, so the link is good only while `r.now()` is before it. `Before` says so, and `!r.now().Before(p.Expires)` catches the instant itself.

Check everything before you change anything: only once the token matches and the time is right does `delete(r.pending, user)` spend the reset.

If the link still works at 09:30:00 exactly, the code compares with `After` instead of `Before`; if a wrong token spends the real link, the `delete` comes before the checks.
