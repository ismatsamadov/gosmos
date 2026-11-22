`LimitError`'s `Error` method has a pointer receiver. What goes in front of the struct literal so that `checkTemp` returns a `*LimitError`?

In `pagerDetail`, call `errors.AsType` with `*LimitError` in square brackets. It returns the error and an `ok` to check before you read the fields.

Format the line with `fmt.Sprintf`, using `%.1f` for the value and both limits, and return `""` when `ok` is `false`.
