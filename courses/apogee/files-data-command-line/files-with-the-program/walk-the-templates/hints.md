Walk from `"templates"`, not from `"."`: the station table sits beside that folder and is not a template.

In the function you hand to `WalkDir`, look at `err` before anything else, since `d` can be nil when it is set. Return `err` to stop the walk; `WalkDir` then returns it to you.

Folders reach your function too. Append the path only when `d.IsDir()` is false, return nil to keep walking, and return the slice with the error `WalkDir` gave back.
