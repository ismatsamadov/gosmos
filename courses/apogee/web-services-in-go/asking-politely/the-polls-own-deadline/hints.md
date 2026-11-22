Only the line that builds the request changes. `Do` already stops a request when its context ends, as long as the request carries the poll's context.

`http.NewRequestWithContext` takes the same arguments as `http.NewRequest`, with the context first.

If the Cape Wren test still waits a whole second, the request isn't carrying `ctx`. If the canceled test gets a reading back, the context was never passed in either.
