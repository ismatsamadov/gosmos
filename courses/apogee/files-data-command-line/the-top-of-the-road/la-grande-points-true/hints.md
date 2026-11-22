Start by planning the day with `planDay`, and return at once if it fails: nothing has been written yet, and nothing should be.

Build the folder's path from `dir`, `"reports"` and `"daily"` with `filepath.Join`, and create it with `os.MkdirAll`. `saveDay` takes a function of one `io.Writer`, while `writeReport` also needs the plan: write a function literal that takes the writer and passes it to `writeReport` with `plan`. The file's name is the plan's `Day` followed by `.txt`.

Write the report to `stdout` with `writeReport` only after `saveDay` has succeeded, and return that error too. If the test finds new files in the current folder, a path was built without `dir` in front.
