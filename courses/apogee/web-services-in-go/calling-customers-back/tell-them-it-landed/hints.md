Build the request with `http.NewRequestWithContext`: `ctx`, the POST method, `url`, and a reader over `body` from `bytes.NewReader`. Set its `Content-Type` header before you send it.

Send it with `c.Do`. When that returns an error there is no answer to close; otherwise defer closing `resp.Body` right after the check, so every return after it closes the body.

A 2xx is any status from 200 to 299, both ends included. For any other, return an `&AnswerError{...}` with the answer's `StatusCode` and `Status`. If only the 202 and 204 cases fail, your check accepts 200 alone.
