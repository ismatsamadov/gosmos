In `board.html`, the note's cell calls `asHTML` on `.Note`. Print `.Note` on its own, as the other cells print their fields, and the template escapes it like every other value.

In `page.go`, `parsePage` hands the template its functions through `Funcs` and a `template.FuncMap`. With `asHTML` gone, nothing needs that map: `template.ParseFS` parses the two files on its own, and the template it returns is named after the first of them.

Delete the `asHTML` function as well, so no template can call it again. If `NewBoard` fails with `function "asHTML" not defined`, the template still calls it somewhere.
