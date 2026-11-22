## Which layout writes a time as 2027-09-22?

- [ ] YYYY-MM-DD
  > Go layouts don't use codes such as YYYY: Format prints these letters exactly as they are.
- [x] 2006-01-02
  > Year, month and day of the reference time stand for the same parts of your time.
- [ ] 2027-09-22
  > A layout uses the reference time; the digits of your own date are read as layout parts and come out garbled.

## Are 2027-09-22T07:14:00-07:00 and 2027-09-22T14:14:00Z the same instant?

- [x] Yes: 07:14 seven hours behind UTC is 14:14 UTC
  > One moment, written once with the van's offset and once in UTC.
- [ ] No, they are seven hours apart
  > The offset is part of the timestamp; once it is applied, the two readings name the same moment.
- [ ] Only on a laptop set to seven hours behind UTC
  > A complete timestamp carries its own offset, so every laptop reads the same instant from it.

## What does t.UTC() change about t?

- [ ] The instant, moved by the offset
  > The instant stays where it is; t.Equal(t.UTC()) is always true.
- [ ] Nothing that Format can show
  > Formatting t.UTC() writes the UTC clock and a Z instead of the original offset.
- [x] Only its location: the clock it is shown on
  > It returns the same instant set to UTC, ready to be written the network's way.
