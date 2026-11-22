## The search has `uplinkd.log` open when the rotator renames it. What does it read next?

- [x] Whatever is left of the same file, now named `uplinkd-2028-12-07.1.log`
  > A rename changes only the file's name: no byte is copied, and a program holding the file open goes on reading the same file.
- [ ] The fresh `uplinkd.log`, since that is the name it opened
  > An open file stays the file it was; only a tool that opens the name again finds the fresh one.
- [ ] Nothing more: the rename copies the file and removes the original
  > On one disk a rename moves no data at all, and the file the reader holds is still there.

## Why can't a rotation cut a record in two?

- [ ] The rotator waits for a newline before it rotates
  > It never looks inside `p`: it decides before each `Write`, whatever the bytes are.
- [x] A record is one `Write`, and the rotator rotates only between writes
  > The JSON handler writes each record with a single `Write`, so every record lands whole in one file.
- [ ] Records are never longer than `MaxBytes`
  > A longer one goes alone into a fresh file; a record's length has nothing to do with splitting it.

## Why does the current file always keep the name `uplinkd.log`?

- [ ] A file can't be renamed while a program has it open
  > It can, as the first question shows.
- [ ] Its day isn't known until it is rotated
  > The rotator knows the current file's day from the start and keeps it in `r.day`.
- [x] Tools and people always find the newest lines under one name
  > A tool that follows `uplinkd.log` by name, opening it again after each rotation, always reads the file being written.
