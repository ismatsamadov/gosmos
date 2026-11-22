`number` receives a copy of the frame, so it can set `f.Seq` and return `f`. Take the number from `l.next` first, then add one to `l.next` for the command after it.

In `learn`, the controller's next free number is `w.LastSeq + 1`. The built-in `max` returns the larger of two values of the same type, and both are `uint32` here.

If the test says `learn(last_seq 57)` left `l.next` at 58, the link went back to numbers it had used: keep `l.next` when it is already larger.
