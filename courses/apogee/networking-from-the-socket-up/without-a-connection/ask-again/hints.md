Wrap the sending, the deadline and `readAnswer` in a loop that runs `tries` times, so every try sends `question(n)` and gets a fresh deadline of `wait` from now.

After `readAnswer`, check its error with `errors.Is(err, os.ErrDeadlineExceeded)`: only a deadline means try again. A nil error is the answer; any other error ends `ask` at once.

Keep the last error in a variable declared before the loop, so the deadline error from the last try can be returned after it. Every try asks with the same `n`.
