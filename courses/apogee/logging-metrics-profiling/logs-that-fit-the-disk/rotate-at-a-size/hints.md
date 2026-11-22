In `rotate`, close `r.f` first, then move the current file with `os.Rename`, from its path in `r.Dir` (the program's name and `.log`) to the path of the name the task gives. Return what `r.open()` returns: it opens a fresh `uplinkd.log` and counts the day's rotated files again.

In `Write`, the check goes in place of the TODO line, before the write: rotate only when the file already holds bytes and its size plus the record's length is above `r.MaxBytes`. `len(p)` is an `int`, so convert it to add it to `r.size`, and return 0 and the error if `rotate` fails.

If a file is rotated at exactly 1024 bytes, or an empty file is rotated, compare your check with the task's two conditions: a file may reach the limit exactly, and an empty file takes any record, however long.
