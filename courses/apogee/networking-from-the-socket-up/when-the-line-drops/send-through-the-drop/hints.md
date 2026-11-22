Number the command once, before a loop, and keep a pointer to it in `l.inflight`. Inside the loop, each try is one call to `l.client.Do` with that same frame.

An answer with no error ends the loop: clear `l.inflight` and return what `settle` says. An error that `lost` doesn't recognize, such as the context's, goes straight back to the caller.

On a lost line, close `l.client` and call `l.connect(ctx)`, returning its error if it fails; then a nil `l.inflight` means the command ran, so return nil. If the test shows the STOW executed as 121, the command was numbered again inside the loop.
