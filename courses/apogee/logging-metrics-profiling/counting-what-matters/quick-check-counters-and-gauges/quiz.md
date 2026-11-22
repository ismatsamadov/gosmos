## `passes_total` fell from 41 to 3 between the readings. What happened?

- [ ] 38 passes were taken back after they were counted
  > A counter only ever goes up: nothing in uplinkd takes a pass back.
- [x] uplinkd restarted, and has counted 3 passes since
  > A counter starts again at 0 when its program restarts, so a drop means a restart, and the new value is what it counted since.
- [ ] The 16:00 reading is wrong and should be thrown away
  > The reading is right: it is the new program's count, and a total over the day still needs it.

## Why is `passes_running` a gauge and not a counter?

- [x] It goes down whenever a pass ends
  > A gauge is a reading that moves both ways, and a pass that ends takes its 1 away.
- [ ] It changes more often than a counter does
  > How often a number changes doesn't decide its kind; whether it can go down does.
- [ ] It stays small, and counters are for big numbers
  > Size doesn't decide its kind either: a counter may read 3 and a gauge 20.
