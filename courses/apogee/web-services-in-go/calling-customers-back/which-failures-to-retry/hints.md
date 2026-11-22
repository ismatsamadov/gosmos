Check `errors.Is(err, context.Canceled)` first: it means the service is stopping, so the answer is false whatever the status says.

With `status` 0 no answer came, so retry exactly when there is an error. With any other status, decide by the status alone: `err` is then only the `*AnswerError` that carries it.

If 404 or 410 comes out true, you are retrying on any error at all; look at the status before the error. A `switch` with no value after it, its cases in this order, keeps the rules readable.
