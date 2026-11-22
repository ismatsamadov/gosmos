`http.Get` returns the response and an error. When the error isn't nil there is no response to read, so return it with two empty strings.

Right after that check, defer closing `resp.Body`, then read all of it with `io.ReadAll`, which returns the bytes and an error of its own.

`resp.Status` is the whole status, such as `404 Not Found`, while `resp.StatusCode` is the number alone. Return the body as it came: `main` trims the newline at its end.
