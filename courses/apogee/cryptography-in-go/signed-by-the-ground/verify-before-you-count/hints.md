Move the lines that build the canonical form and call `verify` to the top of `Accept`, so a command that fails them returns before anything else runs.

Reading a nil map is safe and gives 0, so the comparison with the last number can come before the map is made. Make the map only on the way to recording, which is the last thing `Accept` does.

If the forged command still blocks command 119, something writes to `g.last` before `verify` has said yes: every refusal must return with the map exactly as it was.
