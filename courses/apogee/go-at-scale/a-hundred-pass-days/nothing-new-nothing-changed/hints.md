Call `Plan` with a window from `at` to `cur.To` and the hand-over's free times: it plans every pass that rises at `at` or later, from the dishes `cur`'s earlier contacts leave.

`slices.IndexFunc` finds the first of `cur.Contacts` whose AOS isn't before `at`; the contacts before it are the ones to keep. A result of -1 means every contact rose before `at`.

Build the result's contacts with `slices.Concat`, then give it `cur.Window`. If the 07:48:33 test finds one contact too many, a contact that rises at `at` is both kept and planned again.
