In `writeEvent`, marshal the view with `json.Marshal`, then write the three lines and the blank line in one `fmt.Fprintf`: `id:` takes `v.version`, and `data:` takes the JSON as a string, all on its one line. Return the first error you meet.

The event is still in the server's buffer after `Fprintf`. `http.NewResponseController(w).Flush()` sends it to the screen and returns an error if it can't; return that error too.

Keep the two given lines at the top of `serveEvents`: they count the stream while it is open. After the headers and the first event, a receive from `r.Context().Done()` waits until the screen goes away. Add `"encoding/json"` to the imports.
