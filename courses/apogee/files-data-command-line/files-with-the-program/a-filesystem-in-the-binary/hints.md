The directive goes on its own line between the doc comment and `var bundled embed.FS`, and names the folder as it is: `templates`, with no slash and no `./`.

Inside `bundled`, every name starts with the folder it was embedded from. Build that name from `"templates"` and `name` with `path.Join`, and pass it to `bundled.ReadFile`.

Return the error from `ReadFile` unchanged: for a template that isn't there it already reports `fs.ErrNotExist`. On success, convert the bytes with `string`.
