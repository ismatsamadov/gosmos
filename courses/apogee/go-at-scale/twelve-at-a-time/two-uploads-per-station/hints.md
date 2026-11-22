Lock `u.mu`, look the station up in `u.byStation`, and when it isn't there yet, make it with `newSlots(u.per)` and store it. Unlock before anything waits.

Then call `Acquire` on the station's slots with `ctx`, returning its error if it has one, and `defer` the `Release` so the slot comes back however `upload` ends. Return what `upload` returns.

If the test reports two uploads in flight where it wants four, every station shares one `Slots`. If an upload runs after its context was canceled, `upload` is called without waiting for a slot first.
