Open the zip with `zip.OpenReader(path)` and return its error at once: a missing file or one that isn't a zip fails right there. Once it is open, `defer` its `Close`.

Loop over `zr.File` and skip every file whose `Name` is not `manifest.xml`. The packer decides the order, so the manifest need not come first.

For the manifest, call its `Open`, close that reader when you are done, read it whole with `io.ReadAll` and return what `readManifest` returns. If the loop ends without finding it, return an error.
