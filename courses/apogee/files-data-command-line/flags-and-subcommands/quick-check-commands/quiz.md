## Where should gctl passes write its booking lines?

- [ ] To standard error, so they stand out
  > Standard error is for what went wrong; a script that wants the bookings reads standard output.
- [ ] To both, so nothing is missed
  > A copy on standard error would mix results into the diagnostics.
- [x] To standard output
  > Results go to standard output, where people and scripts read them.

## Where does the line "gctl: open week.json: no such file or directory" belong?

- [ ] On standard output, after any results
  > A script reading the results would take it for one more booking line.
- [x] On standard error
  > It says what went wrong with the run, so it stays out of the output scripts read.
- [ ] Nowhere, since the status already says it failed
  > The status tells a script that it failed; the line tells a person why.

## Why does run return the status instead of calling os.Exit itself?

- [x] So a test can call run and check the status
  > os.Exit would end the whole test process; a returned status can be compared like any other value.
- [ ] Because os.Exit only works inside main
  > os.Exit works anywhere, which is exactly why calling it in run would end a test.
- [ ] Because returning a number is faster than exiting
  > Speed has nothing to do with it: main still exits with the same status in the end.
