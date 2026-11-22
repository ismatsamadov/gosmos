What should a test do when its input file is missing: carry on with zero lines, or stop right there?

`filepath.Join` builds the path from `"testdata"` and `name`. The error from `os.ReadFile` already names the file, so put that error in your `t.Fatalf` message.

`strings.Split` on `"\n"` gives you the lines, plus an empty string for the blank line between the passes and one after the final newline: keep only the others.
