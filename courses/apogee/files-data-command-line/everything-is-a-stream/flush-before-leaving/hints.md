Send the formatted lines to the buffer rather than directly to `w`. Each write can still fail before the loop finishes.

After the loop succeeds, return the result of `Flush`. Otherwise a small buffered report may never reach the caller's writer.
