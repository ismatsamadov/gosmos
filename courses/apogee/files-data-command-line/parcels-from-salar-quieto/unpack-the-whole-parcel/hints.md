`zip.OpenReader` gives you the parcel's file list in `zr.File`. Hand the whole list to `checkNames` before you open the root, so a single unsafe name stops the parcel with nothing written.

Open `dir` with `os.OpenRoot`, and close the root and the zip when you return. In the loop, give every file to `unpackFile` with `maxFileBytes` as its limit, and return the first error it reports.

After the loop the manifest is one of the unpacked files: read it with the root's `ReadFile`, then let `readManifest` check it and return what it returns.
