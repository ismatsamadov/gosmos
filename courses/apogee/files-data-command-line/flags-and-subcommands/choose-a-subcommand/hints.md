Add a case for each command to the switch. Each one returns what its run function returns, called with the arguments that come after the command's name.

After the switch, every name that did not match is unknown. `fmt.Fprintf` with `%q` writes the name in quotes; follow it with a newline and `usage`, and return 2.

When `stations` answers "stations takes no arguments", the command received its own name too: slice the arguments so they start after it.
