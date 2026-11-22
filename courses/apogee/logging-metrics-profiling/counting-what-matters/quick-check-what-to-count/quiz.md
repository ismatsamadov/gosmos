## Lantern Atoll's log keeps only some `carrier lost` warnings. Why count events with counters, not log lines?

- [ ] Counting lines takes longer than reading one number
  > Speed isn't the trouble: the count itself would be wrong.
- [x] The sampler drops repeated lines, so a count of lines comes out short
  > After the first ten of a minute it writes one in a hundred, while a counter counts every event, written or not.
- [ ] Log lines carry no time, so they can't be counted by day
  > Every line carries its time; what's missing is the lines the sampler never wrote.

## Yusra would like `passes_total` by pass ID too. Why does Nkem say no?

- [ ] A pass ID holds slashes, which a JSON key can't carry
  > A JSON key can hold any text, slashes included.
- [ ] A pass ID is too long to be a map key
  > A Go map key can be a string of any length; the trouble is how many keys there would be.
- [x] Each pass would get a counter of its own, used once and then kept
  > A label's values must be few, such as satellites or outcomes: pass IDs would add about fifteen counters a day at Lantern Atoll alone, each kept for as long as uplinkd runs.

## What can `ack_wait_seconds` answer that an average wait can't?

- [x] How many acknowledgments took longer than 10 seconds
  > The buckets above 10 s count exactly those, while an average of 2 s can hide a five-minute wait.
- [ ] What the average wait was this week
  > The histogram answers that too, as `sum_ms` divided by `count`, and so does the average itself.
- [ ] Which pass waited longest, and when
  > Neither can: a histogram keeps a count per bucket, not the passes behind them.
