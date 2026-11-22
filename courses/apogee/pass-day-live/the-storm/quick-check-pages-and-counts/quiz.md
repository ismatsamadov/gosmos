## In the wind lock, Ystmark skips nine booked Furrows. What do the pager and the week's report make of them?

- [ ] A page each, and nine lost bookings
  > The pass rules count failed passes only, and a page goes out once per problem, never once a pass.
- [x] No page, and nine lost bookings
  > A skipped pass is no failure, so the rules never count it; but a booking whose pass didn't fly is lost, whatever the reason.
- [ ] No page, and nothing lost, since the station chose to stow
  > Lost counts every booked pass that wasn't flown: unplanned, failed or skipped alike.

## Ystmark has 210 bookings and a 99.5% target. How many passes may it lose?

- [ ] 1.05
  > The budget is counted in whole passes, and a station can't lose a twentieth of one.
- [ ] 2
  > Losing 2 of 210 leaves 99.05% flown, under the 99.5% target.
- [x] 1
  > 210 × (10000 − 9950) / 10000 is 1 in whole numbers: losing a second pass would break the target.

## Where should the week's report get its numbers?

- [x] From what the rehearsal counted and kept as the week ran
  > Counters and each pass's outcome are exact the moment the week is done, and need no parsing.
- [ ] From the logs, read back after the week
  > The logs tell people what happened; counting their records back ties the report to their wording, when the rehearsal already holds every count.
- [ ] From the console's lines, read back after the week
  > The console is `w.Out`, an `io.Writer`: the rehearsal writes its lines there and can't read them back.
