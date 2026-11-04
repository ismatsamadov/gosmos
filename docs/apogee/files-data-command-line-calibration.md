# Files, Data and the Command Line: calibration chapter

This is the first chapter, authored and imported before the rest of Season 2. The complete lesson assets live in `courses/apogee.sqlite`; export `files-data-command-line` / `a-schedule-on-a-stick` with the author tool to inspect starters, read-only tests, solutions, hints and quiz answers.

## Calibration decisions

All seven code solutions passed 30 runs; all six test lessons also passed with `-race`.
Library checking and teaching verification pass for the imported chapter. All nine lesson
views were opened and inspected in both light and dark appearances with temporary learner
data. Independent instructional review remains pending.

- Filesystem lessons use tests and `t.TempDir`; only the path-construction lesson is checked by printed output.
- Each helper does one operation and propagates its error. The final lesson composes the supplied folder and replacement helpers.
- Paths used as day labels are program-controlled; arbitrary path validation is outside this chapter.
- Permission-bit checks describe their scope, and complete-file replacement makes no power-loss durability claim.
- Standard library only, Go 1.27. No test requires a memory stick, a connection or real learner files.

## The lessons

### Read the Schedule

Paz hands you a memory stick with the van's next passes. Beyond Port Ellery there will be no signal, so the schedule must come from the file, not a service the laptop might not reach.

`os.ReadFile(path)` reads the whole file and returns its bytes and an error. Check the error before using the bytes. `string(data)` turns a byte slice into text without trimming spaces or line endings:

```go
data, err := os.ReadFile("notes.txt")
if err != nil {
    return err
}
fmt.Println(len(data))
```

## Your task

Implement `readSchedule`. Return the file's exact text on success, or an empty string and the original error when reading fails. Remove the import placeholder. Tests create their own temporary files; no real memory stick is needed.

### Save the Report

Ingrid needs a report the van can bring back even if the road stays offline. An empty file must mean an empty report, not a save that silently did nothing.

`os.WriteFile` writes bytes and returns an error. It creates a missing file or replaces the contents of an existing one. Its permission argument applies when creating a file; it does not change an existing file's permissions:

```go
err := os.WriteFile("note.txt", []byte("Ready\n"), 0600)
if err != nil {
    return err
}
```

`0600` permits the owner to read and write, with no access for others.

## Your task

Implement `saveReport` using `os.WriteFile`, converting `text` to bytes and using creation mode `0600`. Return its error. The parent folder must already exist; a later lesson creates folders.

### Quick Check: Files

The first reports are ready for the van. Nkem checks the two file helpers before anyone depends on them halfway up the road.

Reading and writing each return an error; the caller must inspect it before announcing success. A file that exists can still be empty. Use the helpers you just completed to reason about the three questions.

## Your task

Choose one answer for each question about missing files, byte slices and creation permissions. The checks are about file operations, not the contents of a particular report.

### Paths That Fit

Paz wants one folder per day on the stick. Hand-written slashes have already left a double separator in one path and a leftover `..` in another.

`filepath.Join` combines local path components and cleans repeated separators, `.` and parent steps such as `..`. It builds a path string; it does not create a folder or check that anything exists:

```go
path := filepath.Join("logs", "old", "..", "current.txt")
fmt.Println(path) // logs/current.txt on this Mac
```

## Your task

Replace the two placeholders with `filepath.Join` calls. First join `plans`, `2027-07-12` and `passes.txt`. Then join `plans`, `old`, `..` and `report.txt`. Keep one printed path per line and remove the import placeholder.

### Make Room for Reports

The van's laptop has no reports folder on its first run. Paz needs the tool to prepare it, and to keep working when the same folder is already there tomorrow.

`os.MkdirAll` creates a directory and any missing parents. An existing directory is fine; a regular file blocking the path is an error. The mode supplies creation permissions, subject to the operating system's permission mask:

```go
err := os.MkdirAll(filepath.Join("cache", "images"), 0700)
if err != nil {
    return err
}
```

## Your task

Implement `prepareReports`. Join `root`, `reports` and `daily`, then create that path with mode `0700`. Return the error so the caller cannot mistake a blocked path for a ready folder.

### Who Can Read It

Mirela checks the van's reports before departure. A write can replace a file's contents while keeping its old permissions, so a successful save alone does not prove who can read it.

`os.Stat` returns information about a path. `info.Mode().Perm()` selects its nine permission bits: owner, group and others. A bitwise `&` keeps only the bits selected by a mask:

```go
mode := info.Mode().Perm()
groupBits := mode & 0070
fmt.Printf("group permissions: %o\n", groupBits)
```

## Your task

Implement `ownerOnly`: inspect the path and return whether its group and other bits, selected by `0077`, are all zero. Return a Stat error rather than treating a missing file as private. This checks permission bits, not every possible operating-system access rule.

### Quick Check: Paths

The van's report folders are ready, and Mirela has checked their permission bits. Nkem asks one last set of questions before the code starts replacing a live schedule.

Building a path, creating its folders and inspecting a file are separate operations. The helpers should say which one failed instead of returning success for all three.

## Your task

Choose one answer per question. Use the operations from the preceding lessons: joining a path does no disk work, existing directories are acceptable, and permission bits are only part of a file's mode.

### Replace It Whole

A partial save must not replace the van's last usable schedule. Paz wants readers to see the old file until the new one is complete.

Create a temporary file in the destination's directory, write it, close it and rename it over the destination. Keeping both files in one directory avoids a rename across filesystems. Remove the temporary name on every failure. This pattern controls what readers see; it does not promise survival after a power failure.

```go
file, err := os.CreateTemp("reports", ".draft-*")
if err != nil {
    return err
}
defer os.Remove(file.Name())
defer file.Close()
```

## Your task

Implement `replaceSchedule` with `CreateTemp`, `Write`, `Close` and `Rename`, checking each error. Use `filepath.Dir(path)` for the temporary file's directory. Close it before renaming, and arrange cleanup even when writing or renaming fails.

### The Stick Is Ready

Paz is ready to leave Port Ellery. Before the memory stick goes in the van, one function must prepare the report folder and save the complete day's schedule, with a useful error if either step fails.

Composition keeps each operation's job small: use the folder helper, construct the filename, then call the replacement helper. Return immediately when a prerequisite fails; a caller must never hear success after a failed setup.

```go
if err := prepareCache(root); err != nil {
    return err
}
return saveIndex(filepath.Join(root, "index.txt"))
```

## Your task

Complete `exportDay` using the supplied read-only `prepareReports` and `replaceSchedule` helpers. Save `text` at `root/reports/daily/<day>.txt`. Here `day` is a trusted label supplied by the program, not an arbitrary path from a user. Propagate setup and save errors.
