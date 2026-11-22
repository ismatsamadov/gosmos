---
slug: three-useful-commands
title: Three Useful Commands
type: code_tests
position: 8
difficulty: 7
topics: errors, programs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz wants the van's daily jobs under one command before the road gets steeper. `passes` and `stations` work; recording that a dish was parked out of the wind, a stow, still lives in her old script.

Every gctl command has the same parts: a flag set from `newFlagSet`, `parseFlags` for help and flag mistakes, `usageError` for a wrong number of arguments, and `report` for the result or the failure. `findStation` looks up a code. A command starts like this:

```go
fs := newFlagSet("point", "point CODE")
status, ok := parseFlags(fs, args, stdout, stderr)
if !ok {
	return status
}
```

## Your task

Implement `runStow`, whose usage starts with this synopsis:

```text
stow [-reason TEXT] CODE
```

Without exactly one code after the flags, return a usage error; for a code not in the station table, fail as `passes` does for an unknown `-station`. Otherwise hand `report` one line, as here with and without `-reason`:

```text
CWR stowed: wind 75 km/h
SQT stowed: no reason given
```
