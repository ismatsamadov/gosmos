Pull from both with `iter.Pull` and defer both stops, then pull a first contact from each, keeping each `ok`. Loop while either schedule has one left.

The earliest AOS is the old head's, unless the old schedule has ended or the new head rises first. Collect with two loops like the example's, one per schedule, each appending its head while the head has that AOS and pulling the next.

If Compare was still running after a second, it took its AOS from a schedule that had ended, whose head is the zero `Contact`; after taking 64 MB, a loop kept appending one head without pulling the next. A change at 05:00:11 where none belongs means contacts were matched one pair at a time, not a second at a time.
