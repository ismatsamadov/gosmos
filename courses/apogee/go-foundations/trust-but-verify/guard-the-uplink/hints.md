Sign `labCommand` with `[]byte(labKey)` first. Each forgery reuses that command, signature and key, changing exactly one of them.

A table with a name, command, signature, key and wanted answer for each case, run with `t.Run(c.name, ...)`, keeps `checkVerify` short.

If a mutant survives, read it in `main_test.go` and work out which case would catch it; the mutant that refuses everything fails only on the valid command. For the forged signature, keep the first 63 digits and end with a different one.
