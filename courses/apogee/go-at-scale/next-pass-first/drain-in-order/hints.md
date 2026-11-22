`v, ok := q.Pop()` takes the next value off the queue: loop on it inside the function `Drain` returns, and when `ok` is false the queue is empty and the function returns.

Give each value to `yield(v)`, and return as soon as it returns false: the loop in `main` has stopped after five passes.

If the program panics with "continued iteration", `Drain` calls `yield` again after it returned false. If it prints nothing, the stub's empty function is still there.
