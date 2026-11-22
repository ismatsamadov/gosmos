## The fourth of six bookings names station SQX. When should gctl find out?

- [ ] While printing the fourth line of the plan
  > Three lines would already be out, and whoever reads them sees a plan that looks complete.
- [x] Before it writes any line of the plan
  > Checking the whole batch first means a failure leaves no partial plan anywhere.
- [ ] After saving, by reading the report back
  > By then a plan made from a bad batch has already replaced last night's on the stick.

## gctl day cannot read settings.json. Where does the explanation go?

- [ ] Standard output, with exit status 0
  > Standard output is for results, and status 0 tells a script that everything worked.
- [ ] Standard error, with exit status 2
  > Status 2 means the command was used wrongly; a file it cannot read is an operational failure.
- [x] Standard error, as one gctl: line, with exit status 1
  > Diagnostics go to standard error, and 1 marks an operational failure.

## What must happen before the temporary report is renamed over last night's plan?

- [x] Every line written and the file closed, both without an error
  > The rename is the only step a reader can see, so it waits until the new report is complete.
- [ ] Last night's plan deleted, to make room
  > Deleting first leaves a moment, or after a failure a whole morning, with no plan at all.
- [ ] Nothing, since the rename checks the contents
  > A rename moves a name; it never looks inside the file.
