`reasons[code]` is the reason phrase for a code, such as `reasons[404]`, which is `Not Found`.

`fmt.Sprintf` with `%d` for the code and `%s` for the reason builds the line; it starts with `HTTP/1.1 ` and ends with `\r\n`, a carriage return and a newline.

If every line still says `200 OK`, the function returns the fixed string. The quotes in the output come from `main`'s `%q`, which also shows the carriage return and newline at your line's end as `\r\n`.
