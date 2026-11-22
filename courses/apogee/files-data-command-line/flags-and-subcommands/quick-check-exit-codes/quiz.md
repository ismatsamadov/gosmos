## What does run return for gctl passes -h?

- [ ] 2, because -h is not one of the command's flags
  > The flag package treats -h as a request for help, and help that was asked for is a success.
- [x] 0
  > The usage goes to standard output, and the command did exactly what was asked.
- [ ] 1, because no bookings were listed
  > Nothing failed: the usage was the whole request.

## What does run return for gctl passes -stn SQT week.json?

- [ ] 1, because no station SQT was found
  > No station was looked up: parsing stopped at the flag the command doesn't define.
- [ ] 0, because week.json is there
  > The file was never opened, since the command line itself was wrong.
- [x] 2
  > A flag the command doesn't define is a usage mistake: the problem and the usage go to standard error.

## What does run return for gctl passes missing.json when that file doesn't exist?

- [x] 1
  > The command line was fine and gctl tried; opening the file failed, so it writes one gctl line to standard error.
- [ ] 2, because the file name is wrong
  > The arguments had the right shape; a file that can't be opened is a failure while working, not a usage mistake.
- [ ] 0, with an empty report
  > An empty report would tell a script there are no passes, which nobody knows yet.
