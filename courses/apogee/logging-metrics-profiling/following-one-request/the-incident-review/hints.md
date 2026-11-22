A root is a span with no `ParentID`, and its children are the spans whose `ParentID` is the root's `SpanID`. Sort a copy of all the spans by `Start` with `slices.SortStableFunc` first, and both the roots and each root's children come out in the order they started.

Walk the sorted spans; for each root, add its line, then walk the sorted spans again and add a line for each of its children. A span's length is `End.Sub(Start)`, and `%8s` prints that duration as `11m5s`, right-aligned in eight columns.

If the children of Cirrine-1's pass come out as `fly` before `worker`, the sort wasn't stable: both started at 01:16:08. If a pass line's columns don't line up, compare the format with the task's: three spaces before a child's name, one before `pass`.
