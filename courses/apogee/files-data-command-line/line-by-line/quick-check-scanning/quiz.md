## The loop ends after 9,000 lines. What does that prove on its own?

- [ ] The log holds 9,000 lines
  > Scan also returns false when a read fails or a record is too long.
- [x] Nothing yet: only sc.Err() can say why it stopped
  > Err returns nil after a clean end and the failure otherwise, so the draft's return nil hides failures.
- [ ] The stick failed at line 9,001
  > A clean end stops the loop the same way, so the loop alone cannot tell.

## Where does the call to sc.Buffer belong?

- [x] Before the first call to Scan
  > The buffer is set up before scanning starts; once Scan has run, Buffer panics.
- [ ] After the loop, where the draft has it
  > By then scanning has started, so Buffer panics.
- [ ] Inside the loop, before each line is read
  > The loop's first Scan has already run there, so Buffer panics.

## With Buffer in the right place, what happens to a 300 KiB record?

- [ ] It is cut at 256 KiB and scanning goes on
  > The scanner never splits a record to fit; it stops instead.
- [ ] The buffer keeps growing until the record fits
  > It grows only up to the limit given to Buffer, here 256 KiB.
- [x] Scan returns false and Err returns bufio.ErrTooLong
  > A record longer than the buffer's limit ends the scan with that error.
