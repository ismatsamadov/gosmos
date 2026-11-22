There may be nothing in flight: check `l.inflight != nil` before reading its `Seq`, or the nil pointer stops the program.

"At or past" is `>=`: a `LastSeq` of 122 means command 122 ran, and so does a `LastSeq` of 150. In that case set `l.inflight` to nil, then return it.

If the test says the controller refused the jog as a duplicate, `resume` sent back a command that had already run. If it ran twice, the frame came back with a new `Seq`: return `l.inflight` itself, number untouched.
