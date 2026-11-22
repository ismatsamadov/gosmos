## Which function spends the most time in its own code?

- [x] `elevationDeg`, with 81% of the time flat
  > Its 188 s of 232 are spent inside `elevationDeg` itself, more than all the other functions together.
- [ ] `planDay`, with 100% cumulative
  > Its cumulative time counts everything it called; its own code took 0.1 s.
- [ ] `scanEverySecond`, with 86% cumulative
  > Most of that is the `elevationDeg` calls it makes; its own loop took 10.2 s.

## What does `planDay`'s line say?

- [ ] It is the slowest function, so it is the first to rewrite
  > Rewriting its own 0.1 s would gain almost nothing: its 232 s belong to the functions it calls.
- [ ] Its 0.1 s means the profile missed most of its work
  > The profile counted every function that ran; `planDay`'s own code took very little of the time.
- [x] It does almost nothing itself: the time is in what it calls
  > Every call began in `planDay`, so all 232 s count in its `cum`, while its `flat`, its own code, is 0.1 s.

## `elevationDeg` is asked about every second, for every satellite and station. Where is the first gain?

- [ ] Making each call faster, since its code is 81% of the time
  > A faster call saves a little on every one of millions of calls, most of which needn't be made at all.
- [x] Asking it far less often: most of those seconds hold no pass's edge
  > A pass's AOS and LOS are a few seconds in the day, so a planner that asks less often, and closer only near an edge, saves most of the 188 s.
- [ ] Rewriting `planDay`, since every call starts there
  > `planDay`'s own code takes 0.1 s; the gain is in how often its loops ask `elevationDeg`.
