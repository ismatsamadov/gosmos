`Handler` already registers the page with `mux.HandleFunc`. Add a line like it for the pattern `GET /board.json` and the method value `b.serveJSON`, written without parentheses, so the mux calls it for each request.

`Refresh` writes `b.json` while it holds `b.mu`, so read it the same way: lock, copy the slice into a local variable, unlock. The slice itself is never changed after `Refresh` stores it, so the copy is safe to write after the lock is gone.

Set both headers with `w.Header().Set` before writing, then hand the bytes to `w.Write`. Remove the `TODO` line and the `"fmt"` import it needed.
