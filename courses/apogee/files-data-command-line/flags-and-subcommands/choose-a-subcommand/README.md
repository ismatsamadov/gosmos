---
slug: choose-a-subcommand
title: Choose a Subcommand
type: code_tests
position: 4
difficulty: 6
topics: programs, switch, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz retires two more scripts today: from now on, `gctl passes` lists a file's bookings and `gctl stations` prints the station table, on every laptop the same way.

A command with subcommands reads its first argument to choose what to run, then hands the rest of the arguments to that subcommand, which parses its own flags. So `-station` comes after `passes`, and a name the command doesn't know is a usage mistake.

```go
name, rest := args[0], args[1:]
switch name {
case "point":
	return point(rest)
case "park":
	return park(rest)
}
```

## Your task

Finish `run`: send `passes` and `stations` to `runPasses` and `runStations` with the arguments after the command's name. For any other name, write the problem and then the usage to `stderr`, as here for `gctl fly`, and return 2:

```text
unknown command "fly"
usage: gctl <command> [flags] [arguments]
```
