## A pass rises at 23:54 UTC on 22 September and sets at 00:05 UTC on the 23rd. What is los.Sub(aos)?

- [x] 11m0s
  > Sub measures the time between the two instants, so midnight changes nothing.
- [ ] -23h49m0s
  > That is what subtracting the clock readings alone gives; each instant carries its date too.
- [ ] An error, because the dates differ
  > Sub returns a Duration and no error, whatever days the two instants fall on.

## The van's clock is `time.FixedZone("", -7*60*60)`. What does `t.In` on that zone change about `t`?

- [ ] The instant, moved seven hours earlier
  > That is what t.Add(-7 * time.Hour) does; In keeps the instant, so the result is Equal to t.
- [ ] Nothing, because the zone has no name
  > A fixed zone needs only its offset; with an empty name, the result still shows the van's clock.
- [x] How it is written: the van's clock, with -07:00
  > In keeps the instant and sets the location, so Format writes the van's clock and offset.

## Which date does a time.Duration of 11m0s carry?

- [ ] The date of the AOS it was measured from
  > Sub returns only the length; the AOS keeps its own date.
- [x] None: a duration is only a length of time
  > It counts nanoseconds and nothing else, so it carries no date, offset or location.
- [ ] The date on the laptop that computed it
  > Nothing about the laptop's clock goes into a duration.
