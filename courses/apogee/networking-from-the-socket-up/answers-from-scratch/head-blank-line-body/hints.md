Wrap `w` in one `bufio.NewWriter` and write everything to that, never to `w`: the status line from `statusLine`, then `"Date: "`, the formatted time and `"\r\n"`.

Each field is `f.Name + ": " + f.Value + "\r\n"`, in the order of the slice. After the last one, `"\r\n"` on its own is the blank line that ends the head; the body follows as it is, even when it's empty.

The errors of `WriteString` and `Write` on a `bufio.Writer` are kept, so `Flush`'s error is the one to return. If the test at UTC+2 shows `08:40:36`, `now` was formatted without `UTC()`.
