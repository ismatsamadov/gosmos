Start from the defaults, `now` and 7, and replace each one only when its parameter isn't empty: `q.Get` returns an empty string for a parameter the query doesn't have.

Parse `from` with `time.Parse` and the layout `time.RFC3339`, and `days` with `strconv.Atoi`. A `days` that parses is still wrong when it is below 1 or above 14.

If `from=monday` lists passes, the error from `time.Parse` is dropped somewhere. If `days=0` gets a list, the range check is missing. On a failure, return the zero time, 0 and an error with exactly the task's words.
