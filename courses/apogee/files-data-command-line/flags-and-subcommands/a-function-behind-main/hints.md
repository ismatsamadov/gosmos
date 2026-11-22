Loop over `stations` and write one line per station with `fmt.Fprintln`, whose first argument is the writer to use: here, the `stdout` parameter.

`Fprintln` puts a space between its values and a newline after the last, so the code and the name are all it needs. Return 0 after the loop.

If the test reports an empty stdout while the table shows up elsewhere in the output, the lines went to the process's own output: use `Fprintln` with the `stdout` parameter, not `Println`.
