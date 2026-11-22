`url.Parse` returns an error for an address with a tab or a line break in it, so a `next` it can't parse is never a path on this site.

`strings.HasPrefix` checks the start twice: one `/` must be there, and `//` must not. `strings.Contains` finds a backslash anywhere; in Go it is written `"\\"`, or `` `\` `` between backquotes.

If `https:/apogee-bookings.example` gets through, you checked the parsed URL's host: it is empty there, yet a browser goes to another site.
