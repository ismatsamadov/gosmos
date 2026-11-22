Save `runtime.NumGoroutine()` in a variable before calling `f`, and subtract it from the count after `f` returns: that difference is how many goroutines `f` left.

Then wait, for at most a second after `f` returns: while the difference is above 0 and the second isn't over, sleep for `time.Millisecond` and count again. Stop as soon as it reaches 0, so a function that leaves nothing behind costs a few milliseconds.

If the test's 20 ms goroutines come back as 3, the count was read once, before they ended. If a function that ends goroutines gives a number below 0, return 0 instead: `max` takes the larger of two numbers.
