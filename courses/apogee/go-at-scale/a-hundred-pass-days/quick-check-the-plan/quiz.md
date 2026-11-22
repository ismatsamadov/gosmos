## Why can twelve stations be planned at once and still give the plan of one at a time?

- [ ] The group's slots stop two stations from planning the same pass
  > Slots only limit how many jobs run at once; they don't decide what a job may touch.
- [x] No two stations share a dish, so each station's plan depends on its own passes and dishes alone
  > Nothing one station decides can change another's contacts, so the order the jobs run in can't change the plan.
- [ ] `MergeBy` puts every contact back in order, whatever the stations decided
  > The merge only orders contacts; it couldn't repair a plan in which two stations had taken one dish.

## How many contacts does the last `MergeBy` hold at a time?

- [x] Twenty: the next contact of each station's list
  > It pulls one contact from each list into its queue, yields the earliest, and refills only from the list that one came from.
- [ ] All 2,020 of the day, sorted before it yields the first
  > That is the sort the merge replaced: it never holds more than the next contact of each list.
- [ ] Two: the one it yields and the one after it
  > To know which contact comes next it has to compare the next contact of every list, not only two.

## Halfway through a tick's run, news cancels it. What does the test Board show?

- [ ] The stations that finished before the cancel, and nothing for the others
  > `Plan` merges only after `Wait` returns no error; a canceled run never builds a schedule.
- [ ] An empty schedule until the next run publishes
  > A canceled run publishes nothing at all, empty or not, so the current schedule stays.
- [x] The last schedule published: a canceled run returns its error and publishes nothing
  > `Plan` returns `ctx.Err()` without a schedule, so the publisher keeps what it had until a newer run finishes.
