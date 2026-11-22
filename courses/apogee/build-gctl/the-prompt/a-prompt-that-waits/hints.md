Make one scanner from `in` before the loop, and let `for sc.Scan()` visit every line: the loop ends by itself when the input does, including a last line that has no newline.

The prompt comes before every line and once more after the last one. Writing it once before the loop and again at the end of each turn of the loop covers both, and an empty input still gets its one prompt.

Write the prompt in `s.prompt` to `out` with `fmt.Fprint`, which adds no newline, so the night shift types on the same line. In the loop, hand `sc.Text()` and `out` to `s.dispatch`, which writes the answer.
