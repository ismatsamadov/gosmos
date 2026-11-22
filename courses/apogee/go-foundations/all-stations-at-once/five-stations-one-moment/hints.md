Nothing in `main` waits: the printing loop runs while all five polls are still asleep. Where does the program need to stop until they are done?

A `sync.WaitGroup` has the two methods you need: `Go` starts a function in a goroutine and counts it, and `Wait` blocks until every counted goroutine has returned.

Three changes: a `sync.WaitGroup` variable before the loops, `wg.___` in place of the `go` statement, and `wg.___()` right before the printing loop.
