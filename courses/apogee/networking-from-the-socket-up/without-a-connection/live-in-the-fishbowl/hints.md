Keep one `Tracker` for the whole pass, declared before the loop. A deadline from `receive` ends the loop with `break`; any other error is returned as it is.

`strings.Split(string(d), "\n")` gives a datagram's lines. For each, `seqOf(line)` gives its number and whether it has one, and `t.Accept(seq)` says whether to print it with `fmt.Fprintln(out, line)`.

After the loop, `l.asked++`, then `ask(l.Questions, l.Recorder, l.asked, l.Wait, l.Tries)`. A deadline error from it means `silent`; any other error is returned. Print `lastLine(t, answer)` last.
