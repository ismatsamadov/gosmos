## When does gctl read the stations.csv it embeds?

- [x] When go build makes the program
  > The build copies the file's bytes into the program, so running it needs no file at all.
- [ ] Each time gctl starts, from the folder it runs in
  > That would be reading the disk at run time, which is what embedding avoids.
- [ ] The first time stationTable is called
  > The variable already holds the table before any function runs.

## Paz edits templates/day.txt on the laptop after building gctl. Which day template does that gctl use?

- [ ] The edited text
  > The program holds the copy made when it was built; it never reads the file on the disk.
- [x] The text it was built with, until gctl is built again
  > Embedded files change only when the program is rebuilt.
- [ ] An error, because the file changed
  > Nothing compares the disk with the program; the embedded copy is still complete.

## Which name opens the stow notice in bundled?

- [ ] /templates/notices/stow.txt
  > Names in a filesystem are relative: no fs.FS holds a name that starts with a slash.
- [ ] ./templates/notices/stow.txt
  > A name may not contain . or .. as a step, so ./ makes it invalid.
- [x] templates/notices/stow.txt
  > Relative, separated by slashes, and starting with the folder the directive named.
