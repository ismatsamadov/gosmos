Declare a `sync.WaitGroup` before the loop, and move the loop's body into a function literal that you pass to its `Go` method.

Call `Wait` on the group after the loop and before the check for errors: until it returns, some results may not be in place yet.

If the simulators gave up waiting, one dish was still asked only after the other answered, so the `status` call must run inside the goroutine. If the order is wrong, store each result at its index instead of appending it.
