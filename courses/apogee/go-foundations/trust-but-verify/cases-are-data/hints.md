What does `canonical` put between the fields, and where do the arguments go? Work out one case on paper before you type it.

Each case is a `canonicalCase` literal: a name, a `Command{...}` with only the fields you need, and the exact string you expect. A command with no arguments, or one with `Sig` set, makes a good case too.

A skeleton for one case: `{name: "...", cmd: Command{SatID: ..., Seq: ..., Op: ...}, want: "..."},` with a comma after each case.
