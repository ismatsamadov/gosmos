## When does one validator that reads tags pay off over a `Validate` method on each struct?

- [ ] Always: reflection catches mistakes the compiler can't
  > It's the other way round: the compiler checks every comparison in a method when the code builds, and never reads a tag.
- [ ] Never: a method on each struct is always the better choice
  > Each settings struct would repeat the same range checks by hand, and the copies drift apart as fields are added.
- [x] When many structs share the same kinds of rules
  > Here passplan 2's settings, tlmd's and the hosts' forms all do. A method on one type is checked by the compiler and easier to read, so reflection pays only when one validator replaces many methods.

## Someone types `check:"min=1,max=6A"` on `Slots`. When does anyone find out?

- [x] When `Validate` first checks `Slots` and reports `rule max=6A does not fit int`
  > The compiler never reads inside a tag, so its mistakes show up only when the code runs; a test that validates every settings struct finds them before a deploy does.
- [ ] When the program builds: the compiler checks every tag
  > The compiler keeps a tag as text, and `go vet` checks only its `key:"value"` form, which this one has.
- [ ] Never: `checkRule` skips a bound it can't read
  > `measure` returns false for it, and `checkRule` reports the rule as a misfit.

## Should the planner run `Validate` on its settings before every contact it plans?

- [ ] Yes: once compiled, reflection costs what a comparison costs
  > Nkem's timing says otherwise: reading types, tags and fields as the program runs took about 1 µs a `Config`, five hundred times the method's 2 ns.
- [x] No: reflection is far slower than a method, so the check runs once, at start-up
  > A start-up check runs once; in the loop it would run again for every contact, though the settings can't change during a run.
- [ ] Yes, if each check runs in a goroutine of its own
  > A goroutine moves the work, not its cost, and adds its own.
