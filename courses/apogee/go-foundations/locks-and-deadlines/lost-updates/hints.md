`Get` already reads `c.frames` safely. Which of its lines keep other goroutines out, and does `Add` need them too?

After `c.mu.Lock()`, any other goroutine that calls `c.mu.Lock()` waits until `c.mu.Unlock()`, and `defer` makes the unlock happen when `Add` returns. In between, `++` on a map entry counts up from 0 for a new satellite.
