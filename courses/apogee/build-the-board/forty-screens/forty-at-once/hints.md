Move the two lines that subscribe and unsubscribe to the top of `serveEvents`, before the headers and before `b.current()`. Everything published from then on lands in this stream's channel.

Keep `b.current()` in a variable, write it, and remember its `version` in another variable, such as `last`. The version field is unexported, but `serveEvents` is in the same package and can read it.

In the loop, after the check for a closed channel, `continue` when the view's version is not greater than `last`; otherwise write it and set `last` to its version.
