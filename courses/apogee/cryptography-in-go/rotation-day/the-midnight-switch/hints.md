Look for the version named `id` in `ring`. If none has that ID, the answer is false at any time.

Once you have it, three things decide: `Revoked` refuses at once; a zero `Retire`, which `Retire.IsZero()` reports, never retires; otherwise `now` must be before `Retire`. `Activate` plays no part in accepting.

If the rehearsal's command is still refused, the code still compares `now` with the next version's `Activate`. If the version is accepted at its retire instant, compare with `now.Before(v.Retire)`, which is false at that instant.
