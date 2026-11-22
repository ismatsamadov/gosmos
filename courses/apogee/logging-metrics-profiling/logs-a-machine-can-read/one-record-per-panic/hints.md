The failure line is already written by the `l.Error` call: give it one more key and value, `"stack"` and the stack, then delete the `debug.PrintStack()` line.

`debug.Stack()` returns a `[]byte`: make it a `string` before it becomes the value of `stack`, or the JSON handler writes it in base64.

If the test finds no `orbitOf` in the stack, the stack was taken somewhere else: take it right there, in the deferred function, after `recover` has returned the value.
