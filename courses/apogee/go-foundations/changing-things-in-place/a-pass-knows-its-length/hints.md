Look at the example: where does the receiver go, and how does the method's body use it? Your method goes outside `main`, where the `TODO` is.

The receiver is `(p Pass)` and the result type is `time.Duration`. Calling `Sub` on one `time.Time` with another gives the time between them.
