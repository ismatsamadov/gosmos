`b.current()` returns the view the Board shows now, and `b.now()` the time for the heading. Put them in a `pageData` literal, with `View` and `Now` as its field names, and give `b.page.Execute` a pointer to a `bytes.Buffer` as its writer.

When `Execute` returns an error, answer with `http.Error`, the text `internal error` and `http.StatusInternalServerError`, then return. The half-written page stays in the buffer, so the screen gets the error and nothing else.

Set the header with `w.Header().Set` before you write the body: headers leave with the body's first byte, and a change after that is lost. Then send the page with the buffer's `WriteTo` method. Add `"bytes"` to the imports of `page.go`, and remove `"fmt"` once `TODO` is gone.
