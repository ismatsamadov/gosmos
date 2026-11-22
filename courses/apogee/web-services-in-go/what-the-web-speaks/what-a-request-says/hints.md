The path is the part of the URL after `127.0.0.1`. Change the address `main` gives `http.NewRequest` so that its path is `/stations`.

`req.Header.Set` takes two strings: the header's name, then its value. Call it on `req` after the error check and before `printRequest`, so the request you print carries the header.

If the first line reads `GET /stations/ HTTP/1.1`, the URL ends in a slash: the path must be `/stations` exactly. If `Accept` is missing, check that you set it before `printRequest` runs.
