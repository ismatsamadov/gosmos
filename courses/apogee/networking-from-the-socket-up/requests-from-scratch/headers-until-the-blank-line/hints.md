Loop on `readLine`: return its error as it is, and return the header when the line is empty, the end of the head.

`strings.Cut(line, ":")` splits at the first colon only, so a value like `127.0.0.1:8080` keeps its own colon. `strings.ToLower` gives the name its stored form, and `strings.Trim(value, " \t")` takes spaces and tabs off both ends.

Appending to `h[name]` works even for a name the map doesn't hold yet: its slice starts out nil. If `Host: 127.0.0.1:8080` comes back as `8080`, the line was split at its last colon.
