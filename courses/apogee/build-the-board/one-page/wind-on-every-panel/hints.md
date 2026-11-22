Put the definition at the end of `board.html`, after the page: `{{define "panel"}}`, the panel's line, then `{{end}}`. Inside it the dot is one `Panel`, so `.Call`, `.Station` and `.Weather` are its fields, and the class attribute is `panel`, a space, then `.Call`.

`.Weather` is a pointer, nil when the station has no reading: `{{if .Weather}}` takes the wind branch only when there is one, and its `{{else}}` writes `no reading`. In the wind branch, `printf "%.0f"` with `.Weather.WindKmh` prints 41.4 as 41; do the same for `.Weather.GustKmh`.

Below the table, open `<ul class="weather">`, then `{{range .Weather}}` with `{{template "panel" .}}` and a line break before its `{{end}}`, so every panel starts a line of its own. The on-call line prints `.OnCall` inside the `<p>`, all on one line.
