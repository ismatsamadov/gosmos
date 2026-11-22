Ask `fs.Stat` about the `templates` folder first and return its error if it has one: `Sub` alone would hand back a view of a folder that isn't there.

Then return both results of `fs.Sub` for the same folder. Its name has no trailing slash and no `./`.

If `day.txt` is still not found through the view, make sure the function returns the view `Sub` made, not `fsys` itself.
